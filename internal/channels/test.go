package channels

import (
	"fmt"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/utils"
	"slices"
	"strings"
	"time"
)

type TestChannel struct {
	*BaseService
}

func init() {
	RegisterOpenAiChannelService(protocol.ChannelTest, NewTestChannelService)
}

func NewTestChannelService(t *models.ChannelAccount) ChannelOpenApi {
	svc := &TestChannel{
		BaseService: NewBaseService(t),
	}
	return svc
}

// getTestResultByAmount 根据金额范围返回测试结果
// 规则：
// - 金额范围 1-300: 成功
// - 金额范围 301-600: 处理中
// - 金额范围 601以上: 失败
func (t *TestChannel) getTestResultByAmount(amount float64) (status string, resCode, resMsg string) {
	switch {
	case amount >= 1 && amount <= 300:
		return protocol.StatusSuccess, string(protocol.Success), "transaction successful"
	case amount >= 301 && amount <= 600:
		return protocol.StatusPending, string(protocol.StatusPending), "transaction is processing"
	default: // 601以上或0以下
		return protocol.StatusFailed, string(protocol.StatusFailed), "transaction failed"
	}
}

// Payin 实现代收支付请求
func (t *TestChannel) Payin(in *ChannelTrxRequest) *protocol.ChannelResult {
	// 创建日志包装器
	logger := protocol.NewChannelLogWrapper(protocol.ChannelTest, t.AccountID, in)
	var result *protocol.ChannelResult
	defer logger.Log(result)

	// 检查支持的币种
	if !t.isSupportedCurrency(in.Transaction.Ccy) {
		result = &protocol.ChannelResult{
			Status:       protocol.StatusFailed,
			ResCode:      protocol.ResCodeFailure,
			ResMsg:       fmt.Sprintf("Unsupported currency: %s", in.Transaction.Ccy),
			ChannelCode:  protocol.ChannelTest,
			ChannelTrxID: t.generateChannelTrxID(),
		}
		return result
	}

	amt := in.Transaction.Amount.InexactFloat64()
	status, resCode, resMsg := t.getTestResultByAmount(amt)

	result = &protocol.ChannelResult{
		Status:       status,
		ResCode:      resCode,
		ResMsg:       resMsg,
		ChannelCode:  protocol.ChannelTest,
		ChannelTrxID: t.generateChannelTrxID(),
	}

	// 如果是成功状态，设置完成时间
	if status == protocol.StatusSuccess {
		result.CompletedAt = utils.TimeNowMilli()
	}

	return result
}

// Payout 实现代付请求
func (t *TestChannel) Payout(in *ChannelTrxRequest) *protocol.ChannelResult {
	// 创建日志包装器
	logger := protocol.NewChannelLogWrapper(protocol.ChannelTest, t.AccountID, in)
	var result *protocol.ChannelResult
	defer logger.Log(result)

	// 检查支持的币种
	if !t.isSupportedCurrency(in.Transaction.Ccy) {
		result = &protocol.ChannelResult{
			Status:       protocol.StatusFailed,
			ResCode:      protocol.ResCodeFailure,
			ResMsg:       fmt.Sprintf("Unsupported currency: %s", in.Transaction.Ccy),
			ChannelCode:  protocol.ChannelTest,
			ChannelTrxID: t.generateChannelTrxID(),
		}
		return result
	}

	amt := in.Transaction.Amount.InexactFloat64()
	status, resCode, resMsg := t.getTestResultByAmount(amt)

	result = &protocol.ChannelResult{
		Status:       status,
		ResCode:      resCode,
		ResMsg:       resMsg,
		ChannelCode:  protocol.ChannelTest,
		ChannelTrxID: t.generateChannelTrxID(),
	}

	// 如果是成功状态，设置完成时间
	if status == protocol.StatusSuccess {
		result.CompletedAt = time.Now().Unix()
	}

	return result
}

// Refund 实现退款请求
func (t *TestChannel) Refund(in *ChannelTrxRequest) *protocol.ChannelResult {
	// 创建日志包装器
	logger := protocol.NewChannelLogWrapper(protocol.ChannelTest, t.AccountID, in)
	var result *protocol.ChannelResult
	defer logger.Log(result)

	// 检查支持的币种
	if !t.isSupportedCurrency(in.Transaction.Ccy) {
		result = &protocol.ChannelResult{
			Status:       protocol.StatusFailed,
			ResCode:      protocol.ResCodeFailure,
			ResMsg:       fmt.Sprintf("Unsupported currency: %s", in.Transaction.Ccy),
			ChannelCode:  protocol.ChannelTest,
			ChannelTrxID: t.generateChannelTrxID(),
		}
		return result
	}

	amt := in.Transaction.Amount.InexactFloat64()
	status, resCode, resMsg := t.getTestResultByAmount(amt)

	result = &protocol.ChannelResult{
		Status:       status,
		ResCode:      resCode,
		ResMsg:       resMsg,
		ChannelCode:  protocol.ChannelTest,
		ChannelTrxID: t.generateChannelTrxID(),
	}

	// 如果是成功状态，设置完成时间
	if status == protocol.StatusSuccess {
		result.CompletedAt = time.Now().Unix()
	}

	return result
}

// Query 实现查询请求
func (t *TestChannel) Query(in *ChannelTrxRequest) *protocol.ChannelResult {
	// 创建日志包装器
	logger := protocol.NewChannelLogWrapper(protocol.ChannelTest, t.AccountID, in)
	var result *protocol.ChannelResult
	defer logger.Log(result)

	// 模拟查询结果，基于渠道交易ID的特征关键词决定状态
	status := protocol.StatusSuccess
	resCode := protocol.CODE_SUCCESS
	resMsg := "transaction found"

	if in.Transaction.GetChannelTrxID() != "" {
		// 基于渠道交易ID中的关键词来决定查询结果
		trxID := in.Transaction.GetChannelTrxID()
		if strings.Contains(trxID, "failed") || strings.Contains(trxID, "fail") || strings.Contains(trxID, "error") {
			status = protocol.StatusFailed
			resCode = protocol.ResCodeFailure
			resMsg = "transaction failed"
		} else if strings.Contains(trxID, "pending") || strings.Contains(trxID, "processing") || strings.Contains(trxID, "wait") {
			status = protocol.StatusPending
			resMsg = "transaction is processing"
		}
	}

	result = &protocol.ChannelResult{
		Status:       status,
		ResCode:      resCode,
		ResMsg:       resMsg,
		ChannelCode:  protocol.ChannelTest,
		ChannelTrxID: in.Transaction.GetChannelTrxID(),
	}

	// 如果是成功状态，设置完成时间
	if status == protocol.StatusSuccess {
		result.CompletedAt = time.Now().Unix()
	}

	return result
}

// isSupportedCurrency 检查是否支持该币种
func (t *TestChannel) isSupportedCurrency(currency string) bool {
	supportedCurrencies := []string{"USD", "INR", "EUR", "GBP", "JPY", "CNY", "SGD", "HKD"}

	return slices.Contains(supportedCurrencies, currency)
}

// generateChannelTrxID 生成渠道交易ID
func (t *TestChannel) generateChannelTrxID() string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("TEST_%s_%d", t.AccountID, timestamp)
}
