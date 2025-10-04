package protocol

import (
	"encoding/json"
	"fmt"
	"inpayos/internal/log"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

// ChannelLog 表示渠道请求日志的结构体
type ChannelLog struct {
	// 关键业务标识ID
	Mid      string `json:"mid"`        // 商户ID
	ReqID    string `json:"req_id"`     // 商户请求ID
	TrxID    string `json:"trx_id"`     // 系统交易ID
	OriTrxID string `json:"ori_trx_id"` // 原始交易ID(退款时)

	// 渠道信息
	Channel        string `json:"channel"`         // 渠道代码
	ChannelAccount string `json:"channel_account"` // 渠道账户
	ChannelTrxID   string `json:"channel_trx_id"`  // 渠道交易ID

	// 状态信息
	Status        string `json:"status"`         // 系统状态
	ChannelStatus string `json:"channel_status"` // 渠道状态
	ResCode       string `json:"res_code"`       // 响应码
	ResMsg        string `json:"res_msg"`        // 响应信息

	// 关键业务字段
	DealID string `json:"deal_id,omitempty"` // 渠道交易流水号
	Link   string `json:"link,omitempty"`    // 支付链接

	// 费用信息
	ChannelFeeCcy    string           `json:"channel_fee_ccy,omitempty"`    // 渠道费用币种
	ChannelFeeAmount *decimal.Decimal `json:"channel_fee_amount,omitempty"` // 渠道费用金额

	// 时间信息
	CreatedAt   int64  `json:"created_at"`             // 创建时间
	CompletedAt int64  `json:"completed_at,omitempty"` // 完成时间
	Duration    string `json:"duration"`               // 耗时描述
	DurationMS  int64  `json:"duration_ms"`            // 耗时毫秒数

	// HTTP请求信息
	RequestURL     string            `json:"request_url"`     // 请求URL
	RequestHeaders map[string]string `json:"request_headers"` // 请求头
	RequestMethod  string            `json:"request_method"`  // 请求方法
	ResponseStatus int               `json:"response_status"` // HTTP响应状态码

	// 完整请求和响应信息
	BizParams       any            `json:"biz_params"`       // 业务请求参数
	ChannelRequest  MapData        `json:"channel_request"`  // 渠道请求参数
	ChannelResponse MapData        `json:"channel_response"` // 渠道响应内容
	Result          *ChannelResult `json:"result"`           // 系统处理结果
}

// ChannelLogWrapper 用于记录渠道服务的请求日志
type ChannelLogWrapper struct {
	ChannelCode     string
	accountID       string
	channelReqID    string
	startTime       time.Time
	params          any     // 原始业务参数
	channelRequest  MapData // 实际发送给渠道的请求
	channelResponse MapData // 渠道原始响应
}

// NewChannelLogWrapper 创建一个新的日志包装器
func NewChannelLogWrapper(channelCode, accountID string, params any) *ChannelLogWrapper {
	reqID := fmt.Sprintf("%s_%s_%d", channelCode, accountID, time.Now().UnixNano())
	return &ChannelLogWrapper{
		ChannelCode:  channelCode,
		accountID:    accountID,
		channelReqID: reqID,
		startTime:    time.Now(),
		params:       params,
	}
}

// SetRequestURL 设置请求URL
func (w *ChannelLogWrapper) SetRequestURL(url string) *ChannelLogWrapper {
	if w.channelRequest == nil {
		w.channelRequest = MapData{}
	}
	w.channelRequest["request_url"] = url
	return w
}

// SetRequestMethod 设置请求方法
func (w *ChannelLogWrapper) SetRequestMethod(method string) *ChannelLogWrapper {
	if w.channelRequest == nil {
		w.channelRequest = MapData{}
	}
	w.channelRequest["request_method"] = method
	return w
}

// SetRequestHeaders 设置请求头
func (w *ChannelLogWrapper) SetRequestHeaders(headers map[string]string) *ChannelLogWrapper {
	if w.channelRequest == nil {
		w.channelRequest = MapData{}
	}
	w.channelRequest["request_headers"] = headers
	return w
}

// SetRequestBody 设置请求体
func (w *ChannelLogWrapper) SetRequestBody(body MapData) *ChannelLogWrapper {
	if w.channelRequest == nil {
		w.channelRequest = MapData{}
	}
	w.channelRequest["request_body"] = body
	return w
}

// SetResponseStatus 设置响应状态码
func (w *ChannelLogWrapper) SetResponseStatus(status int) *ChannelLogWrapper {
	if w.channelResponse == nil {
		w.channelResponse = MapData{}
	}
	w.channelResponse["response_status"] = status
	return w
}

func (w *ChannelLogWrapper) SetResponseBody(body string) *ChannelLogWrapper {
	if w.channelResponse == nil {
		w.channelResponse = MapData{}
	}
	w.channelResponse["response_body"] = body
	return w
}

// SetResponseData 设置响应内容
func (w *ChannelLogWrapper) SetResponseData(body MapData) *ChannelLogWrapper {
	if w.channelResponse == nil {
		w.channelResponse = MapData{}
	}
	w.channelResponse["response_body"] = body
	return w
}

// SetError 设置错误信息
func (w *ChannelLogWrapper) SetError(err error) *ChannelLogWrapper {
	if w.channelResponse == nil {
		w.channelResponse = MapData{}
	}
	if err != nil {
		w.channelResponse["error"] = err.Error()
	}
	return w
}

// Log 记录完整的请求日志，包含请求参数、响应结果、错误信息和执行时间
func (w *ChannelLogWrapper) Log(result *ChannelResult) {
	duration := time.Since(w.startTime)

	// 提取请求参数中的关键字段
	var reqID, trxID, oriTrxID string
	var mid string
	switch params := w.params.(type) {
	case *ChannelPayinRequest:
		reqID = params.ReqID
		trxID = params.TrxID
		mid = params.Mid
	case *ChannelRefundRequest:
		reqID = params.ReqID
		trxID = params.TrxID
		oriTrxID = params.OriTrxID
		mid = params.Mid
	case *ChannelPayoutRequest:
		reqID = params.ReqID
		trxID = params.TrxID
		mid = params.Mid
	}

	// 构建渠道日志结构体
	channelLog := &ChannelLog{
		// 关键业务标识ID
		Mid:      mid,
		ReqID:    reqID,
		TrxID:    trxID,
		OriTrxID: oriTrxID,

		// 渠道信息
		Channel:        w.ChannelCode,
		ChannelAccount: w.accountID,
		ChannelTrxID:   "",

		// 状态信息
		Status:        "",
		ChannelStatus: "",
		ResCode:       "",
		ResMsg:        "",

		// 时间信息
		CreatedAt:  w.startTime.UnixMilli(),
		Duration:   duration.String(),
		DurationMS: duration.Milliseconds(),

		// 完整请求和响应信息
		BizParams:       w.params,
		ChannelRequest:  w.channelRequest,
		ChannelResponse: w.channelResponse,
		Result:          result,
	}

	// 设置 HTTP 相关信息
	if w.channelRequest != nil {
		channelLog.RequestURL = w.channelRequest.Get("request_url")
		channelLog.RequestMethod = w.channelRequest.Get("request_method")
		if headers := w.channelRequest.GetMapData("request_headers"); len(headers) > 0 {
			headersMap := make(map[string]string)
			for k, v := range headers {
				headersMap[k] = fmt.Sprint(v)
			}
			channelLog.RequestHeaders = headersMap
		}
	}

	if w.channelResponse != nil {
		if status := w.channelResponse.Get("response_status"); status != "" {
			statusInt := cast.ToInt(status)
			channelLog.ResponseStatus = statusInt
		}
	}

	// 添加处理结果的关键字段
	if result != nil {
		channelLog.Status = result.Status
		channelLog.ChannelStatus = result.ChannelStatus
		channelLog.ChannelTrxID = result.ChannelTrxID
		channelLog.ResCode = result.ResCode
		channelLog.ResMsg = result.ResMsg

		// 关键业务字段
		channelLog.DealID = result.DealID
		channelLog.Link = result.Link

		// 费用信息
		channelLog.ChannelFeeCcy = result.ChannelFeeCcy
		channelLog.ChannelFeeAmount = result.ChannelFeeAmount

		// 时间信息
		if result.CompletedAt > 0 {
			channelLog.CompletedAt = result.CompletedAt
		}
	}

	// 输出JSON格式日志
	logJSON, _ := json.MarshalIndent(channelLog, "", "  ")

	// 获取渠道特定的logger
	logger := log.GetServiceLogger(w.ChannelCode)

	// 记录到渠道特定的日志文件
	logger.Info(fmt.Sprintf("[ChannelLog] %s", string(logJSON)))

	// 如果有错误，同时记录到主日志文件
	if result != nil && result.ResCode == ResCodeChannelError {
		log.Get().WithFields(logrus.Fields{
			"req_id":         reqID,
			"trx_id":         trxID,
			"channel_trx_id": result.ChannelTrxID,
			"channel":        w.ChannelCode,
			"error":          result.ResMsg,
		}).Error("Channel request failed")
	}
}

func (w *ChannelLogWrapper) RequestWrapper(URL string, header map[string]string, method string, data, body MapData) *ChannelLogWrapper {
	w.SetRequestURL(URL)
	w.SetRequestMethod(method)
	w.SetRequestHeaders(header)
	w.SetRequestBody(data)
	w.SetRequestBody(body)
	return w
}
