package services

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
)

var webhookService *WebhookService

type WebhookService struct{}

func GetWebhookService() *WebhookService {
	if webhookService == nil {
		webhookService = &WebhookService{}
	}
	return webhookService
}

// CreateWebhook 创建Webhook记录
func (s *WebhookService) CreateWebhook(req *protocol.CreateWebhookRequest) (*models.Webhook, error) {
	webhook := &models.Webhook{
		WebhookID: utils.GenerateID("WH"),
	}
	webhook.WebhookValues = &models.WebhookValues{}
	webhook.WebhookValues.SetUserID(req.UserID).
		SetUserType(req.UserType).
		SetTransactionID(req.TransactionID).
		SetBillID(req.BillID).
		SetType(req.Type).
		SetStatus(req.Status).
		SetAmount(req.Amount).
		SetFee(req.Fee).
		SetCurrency(req.Currency).
		SetNotifyURL(req.NotifyURL).
		SetNotifyStatus("pending").
		SetNotifyTimes(0).
		SetMaxRetryTimes(5).
		SetNextNotifyAt(time.Now().UnixMilli())

	if err := models.DB.Create(webhook).Error; err != nil {
		return nil, fmt.Errorf("failed to create webhook: %w", err)
	}

	return webhook, nil
}

// SendNotification 发送通知
func (s *WebhookService) SendNotification(webhookID string) error {
	webhook, err := s.GetWebhook(webhookID)
	if err != nil {
		return err
	}

	if webhook.WebhookValues.GetNotifyStatus() == "success" {
		return nil // 已经成功，无需重发
	}

	if webhook.WebhookValues.GetNotifyTimes() >= webhook.WebhookValues.GetMaxRetryTimes() {
		return fmt.Errorf("max retry times exceeded")
	}

	// 构建通知数据
	notifyData := map[string]interface{}{
		"webhook_id":     webhook.WebhookID,
		"transaction_id": webhook.WebhookValues.GetTransactionID(),
		"bill_id":        webhook.WebhookValues.GetBillID(),
		"type":           webhook.WebhookValues.GetType(),
		"status":         webhook.WebhookValues.GetStatus(),
		"amount":         webhook.WebhookValues.GetAmount().String(),
		"fee":            webhook.WebhookValues.GetFee().String(),
		"currency":       webhook.WebhookValues.GetCurrency(),
		"timestamp":      time.Now().UnixMilli(),
	}

	// 生成签名
	sign := s.generateSignature(notifyData, "secret_key") // 实际应该从配置获取
	notifyData["sign"] = sign

	jsonData, err := json.Marshal(notifyData)
	if err != nil {
		return fmt.Errorf("failed to marshal notify data: %w", err)
	}

	// 发送HTTP请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(webhook.WebhookValues.GetNotifyURL(), "application/json", bytes.NewBuffer(jsonData))

	// 更新通知记录
	updateData := map[string]interface{}{
		"notify_times":   webhook.WebhookValues.GetNotifyTimes() + 1,
		"last_notify_at": time.Now().UnixMilli(),
		"request_body":   string(jsonData),
	}

	if err != nil {
		// 请求失败
		updateData["response_code"] = "000"
		updateData["response_body"] = err.Error()
		updateData["notify_status"] = "failed"
		updateData["next_notify_at"] = s.calculateNextRetryTime(webhook.WebhookValues.GetNotifyTimes() + 1)
	} else {
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		updateData["response_code"] = fmt.Sprintf("%d", resp.StatusCode)
		updateData["response_body"] = string(respBody)

		if resp.StatusCode == 200 {
			updateData["notify_status"] = "success"
		} else {
			updateData["notify_status"] = "failed"
			updateData["next_notify_at"] = s.calculateNextRetryTime(webhook.WebhookValues.GetNotifyTimes() + 1)
		}
	}

	return models.DB.Model(webhook).Updates(updateData).Error
}

// GetWebhook 获取Webhook记录
func (s *WebhookService) GetWebhook(webhookID string) (*models.Webhook, error) {
	var webhook models.Webhook
	err := models.DB.Where("webhook_id = ?", webhookID).First(&webhook).Error
	if err != nil {
		return nil, fmt.Errorf("webhook not found: %w", err)
	}
	return &webhook, nil
}

// GetPendingWebhooks 获取待通知的Webhook记录
func (s *WebhookService) GetPendingWebhooks(limit int) ([]*models.Webhook, error) {
	var webhooks []*models.Webhook
	now := time.Now().UnixMilli()

	err := models.DB.Where("notify_status = ? AND next_notify_at <= ? AND notify_times < max_retry_times",
		"pending", now).
		Or("notify_status = ? AND next_notify_at <= ? AND notify_times < max_retry_times",
			"failed", now).
		Limit(limit).
		Find(&webhooks).Error

	return webhooks, err
}

// ProcessPendingWebhooks 处理待通知的Webhook
func (s *WebhookService) ProcessPendingWebhooks(ctx context.Context) error {
	webhooks, err := s.GetPendingWebhooks(100)
	if err != nil {
		return fmt.Errorf("failed to get pending webhooks: %w", err)
	}

	for _, webhook := range webhooks {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := s.SendNotification(webhook.WebhookID); err != nil {
				// 记录错误但继续处理其他webhook
				fmt.Printf("Failed to send webhook notification %s: %v\n", webhook.WebhookID, err)
			}
		}
	}

	return nil
}

// generateSignature 生成签名
func (s *WebhookService) generateSignature(data map[string]interface{}, secretKey string) string {
	// 简化的签名算法，实际应该更复杂
	jsonStr, _ := json.Marshal(data)
	hash := md5.Sum(append(jsonStr, []byte(secretKey)...))
	return fmt.Sprintf("%x", hash)
}

// calculateNextRetryTime 计算下次重试时间（指数退避）
func (s *WebhookService) calculateNextRetryTime(retryTimes int32) int64 {
	delays := []int{30, 60, 300, 900, 1800} // 30s, 1m, 5m, 15m, 30m

	var delay int
	if int(retryTimes-1) < len(delays) {
		delay = delays[retryTimes-1]
	} else {
		delay = delays[len(delays)-1]
	}

	return time.Now().Add(time.Duration(delay) * time.Second).UnixMilli()
}

// CreateWebhookFromTransaction 从交易记录创建Webhook
func (s *WebhookService) CreateWebhookFromTransaction(transaction interface{}) error {
	var req protocol.CreateWebhookRequest

	// 根据不同的业务类型转换
	switch tx := transaction.(type) {
	case *models.Receipt:
		req = protocol.CreateWebhookRequest{
			UserID:        tx.ReceiptValues.GetUserID(),
			UserType:      tx.ReceiptValues.GetUserType(),
			TransactionID: tx.RecordID,
			BillID:        tx.ReceiptValues.GetBillID(),
			Type:          "receipt",
			Status:        tx.ReceiptValues.GetStatus(),
			Amount:        tx.ReceiptValues.GetAmount(),
			Fee:           tx.ReceiptValues.GetFee(),
			Currency:      tx.ReceiptValues.GetCurrency(),
			NotifyURL:     tx.ReceiptValues.GetNotifyURL(),
		}
	case *models.Payment:
		req = protocol.CreateWebhookRequest{
			UserID:        tx.PaymentValues.GetUserID(),
			UserType:      tx.PaymentValues.GetUserType(),
			TransactionID: tx.RecordID,
			BillID:        tx.PaymentValues.GetBillID(),
			Type:          "payment",
			Status:        tx.PaymentValues.GetStatus(),
			Amount:        tx.PaymentValues.GetAmount(),
			Fee:           tx.PaymentValues.GetFee(),
			Currency:      tx.PaymentValues.GetCurrency(),
			NotifyURL:     tx.PaymentValues.GetNotifyURL(),
		}
	case *models.Refund:
		req = protocol.CreateWebhookRequest{
			UserID:        tx.RefundValues.GetUserID(),
			UserType:      tx.RefundValues.GetUserType(),
			TransactionID: tx.RecordID,
			BillID:        tx.RefundValues.GetBillID(),
			Type:          "refund",
			Status:        tx.RefundValues.GetStatus(),
			Amount:        tx.RefundValues.GetAmount(),
			Fee:           tx.RefundValues.GetFee(),
			Currency:      tx.RefundValues.GetCurrency(),
			NotifyURL:     tx.RefundValues.GetNotifyURL(),
		}
	case *models.Deposit:
		req = protocol.CreateWebhookRequest{
			UserID:        tx.DepositValues.GetUserID(),
			UserType:      tx.DepositValues.GetUserType(),
			TransactionID: tx.RecordID,
			BillID:        tx.DepositValues.GetBillID(),
			Type:          "deposit",
			Status:        tx.DepositValues.GetStatus(),
			Amount:        tx.DepositValues.GetAmount(),
			Fee:           tx.DepositValues.GetFee(),
			Currency:      tx.DepositValues.GetCurrency(),
			NotifyURL:     tx.DepositValues.GetNotifyURL(),
		}
	case *models.Withdraw:
		req = protocol.CreateWebhookRequest{
			UserID:        tx.WithdrawValues.GetUserID(),
			UserType:      tx.WithdrawValues.GetUserType(),
			TransactionID: tx.RecordID,
			BillID:        tx.WithdrawValues.GetBillID(),
			Type:          "withdraw",
			Status:        tx.WithdrawValues.GetStatus(),
			Amount:        tx.WithdrawValues.GetAmount(),
			Fee:           tx.WithdrawValues.GetFee(),
			Currency:      tx.WithdrawValues.GetCurrency(),
			NotifyURL:     tx.WithdrawValues.GetNotifyURL(),
		}
	default:
		return fmt.Errorf("unsupported transaction type")
	}

	// 只有在有通知URL时才创建webhook
	if req.NotifyURL != "" {
		_, err := s.CreateWebhook(&req)
		return err
	}

	return nil
}
