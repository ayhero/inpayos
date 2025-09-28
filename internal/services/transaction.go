package services

import (
	"fmt"

	"gorm.io/gorm"

	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
)

var transactionService *TransactionService

type TransactionService struct{}

func GetTransactionService() *TransactionService {
	if transactionService == nil {
		transactionService = &TransactionService{}
	}
	return transactionService
}

// CreateReceipt 创建代收记录
func (s *TransactionService) CreateReceipt(req *protocol.CreateTransactionRequest) (protocol.ErrorCode, *models.Receipt) {
	var result *models.Receipt
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		receipt := &models.Receipt{
			RecordID: utils.GenerateID("RCP"),
		}
		receipt.ReceiptValues = &models.ReceiptValues{}
		receipt.ReceiptValues.SetUserID(req.UserID).
			SetUserType(req.UserType).
			SetBillID(req.BillID).
			SetStatus("pending").
			SetAmount(req.Amount).
			SetFee(req.Fee).
			SetCurrency(req.Currency).
			SetChannelCode(req.ChannelCode).
			SetPaymentMethod(req.PaymentMethod).
			SetNotifyURL(req.NotifyURL)

		if req.ExpiredAt > 0 {
			receipt.ReceiptValues.SetExpiredAt(req.ExpiredAt)
		}

		if err := tx.Create(receipt).Error; err != nil {
			return fmt.Errorf("failed to create receipt: %w", err)
		}

		// 创建Webhook通知
		webhookService := GetWebhookService()
		if err := webhookService.CreateWebhookFromTransaction(receipt); err != nil {
			// Webhook创建失败不影响主业务，只记录错误
			fmt.Printf("Failed to create webhook for receipt %s: %v\n", receipt.RecordID, err)
		}

		result = receipt
		return nil
	})

	if err != nil {
		return protocol.DatabaseError, nil
	}
	return protocol.Success, result
}

// CreatePayment 创建代付记录
func (s *TransactionService) CreatePayment(req *protocol.CreateTransactionRequest) (*models.Payment, error) {
	var result *models.Payment
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		// 首先冻结商户资金
		accountService := GetAccountService()
		freezeReq := &protocol.UpdateBalanceRequest{
			UserID:       req.UserID,
			UserType:     req.UserType,
			Currency:     req.Currency,
			Operation:    "freeze",
			Amount:       req.Amount,
			BusinessType: "payment",
			Description:  "Payment freeze",
		}
		if err := accountService.UpdateBalance(freezeReq); err != nil {
			return err
		}

		payment := &models.Payment{
			RecordID: utils.GenerateID("PAY"),
		}
		payment.PaymentValues = &models.PaymentValues{}
		payment.PaymentValues.SetUserID(req.UserID).
			SetUserType(req.UserType).
			SetBillID(req.BillID).
			SetStatus("pending").
			SetAmount(req.Amount).
			SetFee(req.Fee).
			SetCurrency(req.Currency).
			SetChannelCode(req.ChannelCode).
			SetPaymentMethod(req.PaymentMethod).
			SetNotifyURL(req.NotifyURL)

		if req.ExpiredAt > 0 {
			payment.PaymentValues.SetExpiredAt(req.ExpiredAt)
		}

		if err := tx.Create(payment).Error; err != nil {
			return fmt.Errorf("failed to create payment: %w", err)
		}

		// 创建Webhook通知
		webhookService := GetWebhookService()
		if err := webhookService.CreateWebhookFromTransaction(payment); err != nil {
			fmt.Printf("Failed to create webhook for payment %s: %v\n", payment.RecordID, err)
		}

		result = payment
		return nil
	})
	return result, err
}

// CreateRefund 创建退款记录
func (s *TransactionService) CreateRefund(req *protocol.CreateTransactionRequest) (*models.Refund, error) {
	if req.SourceTxID == "" {
		return nil, fmt.Errorf("source transaction ID is required for refund")
	}

	var result *models.Refund
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		refund := &models.Refund{
			RecordID: utils.GenerateID("REF"),
		}
		refund.RefundValues = &models.RefundValues{}
		refund.RefundValues.SetUserID(req.UserID).
			SetUserType(req.UserType).
			SetBillID(req.BillID).
			SetOriginalBillID(req.SourceTxID).
			SetStatus("pending").
			SetAmount(req.Amount).
			SetFee(req.Fee).
			SetCurrency(req.Currency).
			SetNotifyURL(req.NotifyURL)

		if err := tx.Create(refund).Error; err != nil {
			return fmt.Errorf("failed to create refund: %w", err)
		}

		// 扣除商户资金
		accountService := GetAccountService()
		debitReq := &protocol.UpdateBalanceRequest{
			UserID:        req.UserID,
			UserType:      req.UserType,
			Currency:      req.Currency,
			Operation:     "subtract",
			Amount:        req.Amount,
			TransactionID: refund.RecordID,
			BusinessType:  "refund",
			Description:   fmt.Sprintf("Refund for transaction %s", req.SourceTxID),
		}
		if err := accountService.UpdateBalance(debitReq); err != nil {
			return err
		}

		// 创建Webhook通知
		webhookService := GetWebhookService()
		if err := webhookService.CreateWebhookFromTransaction(refund); err != nil {
			fmt.Printf("Failed to create webhook for refund %s: %v\n", refund.RecordID, err)
		}

		result = refund
		return nil
	})
	return result, err
}

// CreateDeposit 创建充值记录
func (s *TransactionService) CreateDeposit(req *protocol.CreateTransactionRequest) (*models.Deposit, error) {
	var result *models.Deposit
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		deposit := &models.Deposit{
			RecordID: utils.GenerateID("DEP"),
		}
		deposit.DepositValues = &models.DepositValues{}
		deposit.DepositValues.SetUserID(req.UserID).
			SetUserType(req.UserType).
			SetBillID(req.BillID).
			SetStatus("pending").
			SetAmount(req.Amount).
			SetFee(req.Fee).
			SetCurrency(req.Currency).
			SetChannelCode(req.ChannelCode).
			SetNotifyURL(req.NotifyURL)

		if err := tx.Create(deposit).Error; err != nil {
			return fmt.Errorf("failed to create deposit: %w", err)
		}

		// 创建Webhook通知
		webhookService := GetWebhookService()
		if err := webhookService.CreateWebhookFromTransaction(deposit); err != nil {
			fmt.Printf("Failed to create webhook for deposit %s: %v\n", deposit.RecordID, err)
		}

		result = deposit
		return nil
	})
	return result, err
}

// CreateWithdraw 创建提现记录
func (s *TransactionService) CreateWithdraw(req *protocol.CreateTransactionRequest) (*models.Withdraw, error) {
	var result *models.Withdraw
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		// 首先冻结用户资金
		accountService := GetAccountService()
		freezeReq := &protocol.UpdateBalanceRequest{
			UserID:       req.UserID,
			UserType:     req.UserType,
			Currency:     req.Currency,
			Operation:    "freeze",
			Amount:       req.Amount,
			BusinessType: "withdraw",
			Description:  "Withdraw freeze",
		}
		if err := accountService.UpdateBalance(freezeReq); err != nil {
			return err
		}

		withdraw := &models.Withdraw{
			RecordID: utils.GenerateID("WTH"),
		}
		withdraw.WithdrawValues = &models.WithdrawValues{}
		withdraw.WithdrawValues.SetUserID(req.UserID).
			SetUserType(req.UserType).
			SetBillID(req.BillID).
			SetStatus("pending").
			SetAmount(req.Amount).
			SetFee(req.Fee).
			SetCurrency(req.Currency).
			SetChannelCode(req.ChannelCode).
			SetNotifyURL(req.NotifyURL)

		if err := tx.Create(withdraw).Error; err != nil {
			return fmt.Errorf("failed to create withdraw: %w", err)
		}

		// 创建Webhook通知
		webhookService := GetWebhookService()
		if err := webhookService.CreateWebhookFromTransaction(withdraw); err != nil {
			fmt.Printf("Failed to create webhook for withdraw %s: %v\n", withdraw.RecordID, err)
		}

		result = withdraw
		return nil
	})
	return result, err
}

// ToTransaction 将业务记录转换为统一的Transaction实体
func (s *TransactionService) ToTransaction(businessRecord interface{}) *models.Transaction {
	transaction := &models.Transaction{}
	transaction.TransactionValues = &models.TransactionValues{}

	switch record := businessRecord.(type) {
	case *models.Receipt:
		transaction.TransactionID = record.RecordID
		transaction.TransactionValues.SetUserID(record.ReceiptValues.GetUserID()).
			SetUserType(record.ReceiptValues.GetUserType()).
			SetBillID(record.ReceiptValues.GetBillID()).
			SetType(protocol.TxTypeReceipt).
			SetStatus(record.ReceiptValues.GetStatus()).
			SetAmount(record.ReceiptValues.GetAmount()).
			SetFee(record.ReceiptValues.GetFee()).
			SetCurrency(record.ReceiptValues.GetCurrency()).
			SetChannelCode(record.ReceiptValues.GetChannelCode()).
			SetPaymentMethod(record.ReceiptValues.GetPaymentMethod()).
			SetNotifyURL(record.ReceiptValues.GetNotifyURL()).
			SetExpiredAt(record.ReceiptValues.GetExpiredAt())

	case *models.Payment:
		transaction.TransactionID = record.RecordID
		transaction.TransactionValues.SetUserID(record.PaymentValues.GetUserID()).
			SetUserType(record.PaymentValues.GetUserType()).
			SetBillID(record.PaymentValues.GetBillID()).
			SetType(protocol.TxTypePayment).
			SetStatus(record.PaymentValues.GetStatus()).
			SetAmount(record.PaymentValues.GetAmount()).
			SetFee(record.PaymentValues.GetFee()).
			SetCurrency(record.PaymentValues.GetCurrency()).
			SetChannelCode(record.PaymentValues.GetChannelCode()).
			SetPaymentMethod(record.PaymentValues.GetPaymentMethod()).
			SetNotifyURL(record.PaymentValues.GetNotifyURL()).
			SetExpiredAt(record.PaymentValues.GetExpiredAt())

	case *models.Refund:
		transaction.TransactionID = record.RecordID
		transaction.TransactionValues.SetUserID(record.RefundValues.GetUserID()).
			SetUserType(record.RefundValues.GetUserType()).
			SetBillID(record.RefundValues.GetBillID()).
			SetType(protocol.TxTypeRefund).
			SetStatus(record.RefundValues.GetStatus()).
			SetAmount(record.RefundValues.GetAmount()).
			SetFee(record.RefundValues.GetFee()).
			SetCurrency(record.RefundValues.GetCurrency()).
			SetNotifyURL(record.RefundValues.GetNotifyURL()).
			SetSourceTxID(record.RefundValues.GetOriginalBillID())

	case *models.Deposit:
		transaction.TransactionID = record.RecordID
		transaction.TransactionValues.SetUserID(record.DepositValues.GetUserID()).
			SetUserType(record.DepositValues.GetUserType()).
			SetBillID(record.DepositValues.GetBillID()).
			SetType(protocol.TxTypeDeposit).
			SetStatus(record.DepositValues.GetStatus()).
			SetAmount(record.DepositValues.GetAmount()).
			SetFee(record.DepositValues.GetFee()).
			SetCurrency(record.DepositValues.GetCurrency()).
			SetChannelCode(record.DepositValues.GetChannelCode()).
			SetNotifyURL(record.DepositValues.GetNotifyURL())

	case *models.Withdraw:
		transaction.TransactionID = record.RecordID
		transaction.TransactionValues.SetUserID(record.WithdrawValues.GetUserID()).
			SetUserType(record.WithdrawValues.GetUserType()).
			SetBillID(record.WithdrawValues.GetBillID()).
			SetType(protocol.TxTypeWithdraw).
			SetStatus(record.WithdrawValues.GetStatus()).
			SetAmount(record.WithdrawValues.GetAmount()).
			SetFee(record.WithdrawValues.GetFee()).
			SetCurrency(record.WithdrawValues.GetCurrency()).
			SetChannelCode(record.WithdrawValues.GetChannelCode()).
			SetNotifyURL(record.WithdrawValues.GetNotifyURL())
	}

	return transaction
}
