package channels

import (
	"inpayos/internal/models"
	"inpayos/internal/protocol"
)

type BaseService struct {
	AccountID   string         // 账户ID
	Secret      string         // 账户密钥
	ChannelCode string         // 通道代码
	Settings    map[string]any // 通道配置
}

func NewBaseService(t *models.ChannelAccount) *BaseService {
	svc := &BaseService{
		AccountID:   t.GetAccountID(),
		ChannelCode: t.ChannelCode,
		Settings:    t.GetSettings(),
	}
	return svc
}

func (t BaseService) Payment(in *protocol.ChannelPayinRequest) *protocol.ChannelResult {
	return nil
}

func (t BaseService) Refund(in *protocol.ChannelRefundRequest) *protocol.ChannelResult {
	return nil
}

func (t BaseService) Payout(in *protocol.ChannelPayoutRequest) *protocol.ChannelResult {
	return nil
}

func (t BaseService) Query(in *protocol.ChannelQueryQuest) *protocol.ChannelResult {
	return nil
}
