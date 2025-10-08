package channels

import (
	"inpayos/internal/models"
	"inpayos/internal/protocol"
)

// ChannelOpenApi 支付渠道接口
type ChannelOpenApi interface {
	Payin(in *protocol.ChannelPayinRequest) *protocol.ChannelResult
	Refund(in *protocol.ChannelRefundRequest) *protocol.ChannelResult
	Payout(in *protocol.ChannelPayoutRequest) *protocol.ChannelResult
	Query(in *protocol.ChannelQueryQuest) *protocol.ChannelResult
}

var channel_open_api_service_lib = make(map[string]func(*models.ChannelAccount) ChannelOpenApi)

func RegisterOpenAiChannelService(channel_account string, svc func(*models.ChannelAccount) ChannelOpenApi) {
	channel_open_api_service_lib[channel_account] = svc
}

// Channels 渠道服务映射
var channels = make(map[string]ChannelOpenApi) // key: accountID

func RegisterOpenApiChannelHandlerLib(channel_account string, svc ChannelOpenApi) {
	channels[channel_account] = svc
}

func GetOpenApiChannelService(channel_account string) (svc ChannelOpenApi, ok bool) {
	svc, ok = channels[channel_account]
	return
}

func LoadChannelOpenApiService() {
	// 先清空现有的channels映射
	channels = make(map[string]ChannelOpenApi)

	// 获取所有渠道账户
	accounts := models.GetChannelAccounts()
	for _, account := range accounts {
		// 检查是否有对应的渠道服务创建器
		svc_fn, ok := channel_open_api_service_lib[account.ChannelCode]
		if !ok {
			continue
		}
		// 创建服务并注册
		svc := svc_fn(account)
		RegisterOpenApiChannelHandlerLib(account.GetAccountID(), svc)
	}
}

func IsRequestErr(resp protocol.MapData) bool {
	channel_status, res_code, _ := GetChannelResult(resp)
	return IsRequestErrCode(channel_status, res_code)
}

func IsRedirectUrlEmpty(in protocol.MapData) bool {
	return in.Get("redirect_url") == ""
}
func IsRequestErrCode(channel_status, res_code string) bool {
	if channel_status == protocol.ResCodeChannelError || res_code == protocol.ResCodeChannelError {
		return true
	}
	return false
}

func GetChannelResult(resp protocol.MapData) (channelStatus, resCode, resMsg string) {
	channelStatus = resp.Get("channel_status")
	resCode = resp.Get("res_code")
	resMsg = resp.Get("res_msg")
	return
}

func ComposeChannelRequestResult(result protocol.MapData, err error) map[string]interface{} {
	if len(result) == 0 || err != nil {
		result = protocol.MapData{}
		result["channel_status"] = protocol.ResCodeFailure
		result["res_code"] = protocol.ResCodeChannelError
		if err != nil {
			result["res_msg"] = err.Error()
		} else {
			result["res_msg"] = "request fail"
		}
	}
	return result
}

func Fail(in *protocol.ChannelResult) {
	in.Status = protocol.StatusFailed
}
