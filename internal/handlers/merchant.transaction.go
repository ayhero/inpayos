package handlers

import (
	"inpayos/internal/middleware"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TransactionListRequest 交易列表请求
type TransactionListRequest struct {
	TrxType          string `json:"trx_type" binding:"required"`  // 交易类型：payin, payout
	TrxID            string `json:"trx_id"`                       // 交易ID
	ReqID            string `json:"req_id"`                       // 商户订单号
	TrxMethod        string `json:"trx_method"`                   // 交易方式
	TrxMode          string `json:"trx_mode"`                     // 交易模式
	Status           string `json:"status"`                       // 交易状态
	FlowNo           string `json:"flow_no"`                      // 流水号
	ChannelCode      string `json:"channel_code"`                 // 渠道代码
	ChannelAccount   string `json:"channel_account"`              // 渠道账号
	ChannelGroup     string `json:"channel_group"`                // 渠道组
	ChannelTrxID     string `json:"channel_trx_id"`               // 渠道交易ID
	SettleStatus     string `json:"settle_status"`                // 结算状态
	CreatedAtStart   int64  `json:"created_at_start"`             // 开始时间
	CreatedAtEnd     int64  `json:"created_at_end"`               // 结束时间
	CompletedAtStart int64  `json:"completed_at_start"`           // 交易完成开始时间
	CompletedAtEnd   int64  `json:"completed_at_end"`             // 交易完成结束时间
	SettledAtStart   int64  `json:"settled_at_start"`             // 结算开始时间
	SettledAtEnd     int64  `json:"settled_at_end"`               // 结算结束时间
	Page             int    `json:"page" binding:"min=1"`         // 页码
	Size             int    `json:"size" binding:"min=1,max=100"` // 每页记录数
}

// TransactionDetailRequest 交易详情请求
type TransactionDetailRequest struct {
	TrxID   string `json:"trx_id" binding:"required"`   // 交易ID
	TrxType string `json:"trx_type" binding:"required"` // 交易类型
}

// ListTransactions godoc
// @Summary 获取交易列表
// @Description 分页获取交易列表
// @Tags 交易管理
// @Accept json
// @Produce json
// @Success 200 {object} protocol.Result{data=protocol.PageResult{list=[]protocol.Transaction}}
// @Router /transactions/list [post]
func (t *MerchantAdmin) ListTransactions(c *gin.Context) {
	lang := middleware.GetLanguage(c)

	// 绑定请求参数
	var req TransactionListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// 从上下文获取商户信息
	merchant := middleware.GetMerchantFromContext(c)

	// 构建查询参数
	query := &models.TrxQuery{
		Mid:              merchant.Mid,
		TrxType:          req.TrxType,
		TrxID:            req.TrxID,
		ReqID:            req.ReqID,
		TrxMethod:        req.TrxMethod,
		TrxMode:          req.TrxMode,
		Status:           req.Status,
		FlowNo:           req.FlowNo,
		ChannelCode:      req.ChannelCode,
		ChannelAccount:   req.ChannelAccount,
		ChannelGroup:     req.ChannelGroup,
		ChannelTrxID:     req.ChannelTrxID,
		SettleStatus:     req.SettleStatus,
		CreatedAtStart:   req.CreatedAtStart,
		CreatedAtEnd:     req.CreatedAtEnd,
		CompletedAtStart: req.CompletedAtStart,
		CompletedAtEnd:   req.CompletedAtEnd,
		SettledAtStart:   req.SettledAtStart,
		SettledAtEnd:     req.SettledAtEnd,
		Page:             req.Page,
		Size:             req.Size,
	}

	// 调用服务层查询交易列表
	transactions, total, code := services.GetMerchantTransactionService().ListTransactionByQuery(query)
	if code != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(code, lang))
		return
	}

	// 转换为协议格式
	var list []*protocol.Transaction
	for _, trx := range transactions {
		list = append(list, trx.Protocol())
	}

	// 构建分页结果
	pagination := &protocol.Pagination{
		Page: req.Page,
		Size: req.Size,
	}

	// 返回成功结果
	c.JSON(http.StatusOK, protocol.NewSuccessPageResult(list, total, pagination))
}

// TransactionDetail godoc
// @Summary 获取交易详情
// @Description 获取指定交易ID的详细信息
// @Tags 交易管理
// @Accept json
// @Produce json
// @Success 200 {object} protocol.Result{data=protocol.Transaction}
// @Router /transactions/detail [post]
func (t *MerchantAdmin) TransactionDetail(c *gin.Context) {
	lang := middleware.GetLanguage(c)

	// 绑定请求参数
	var req TransactionDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// 从上下文获取商户信息
	merchant := middleware.GetMerchantFromContext(c)
	if merchant == nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.MerchantNotFound, lang))
		return
	}

	// 构建查询参数
	query := &models.TrxQuery{
		Mid:     merchant.Mid,
		TrxType: req.TrxType,
		TrxID:   req.TrxID,
		Page:    1,
		Size:    1,
	}

	// 调用服务层查询交易详情
	transactions, _, code := services.GetMerchantTransactionService().ListTransactionByQuery(query)
	if code != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(code, lang))
		return
	}

	// 检查是否找到交易
	if len(transactions) == 0 {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.TransactionNotFound, lang))
		return
	}

	// 转换为协议格式并返回
	transactionInfo := transactions[0].Protocol()
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(transactionInfo, lang))
}

// TodayStatsRequest 今日统计请求
type TodayStatsRequest struct {
	TrxType string `json:"trx_type" binding:"required"` // 交易类型：payin, payout
}

// GetTransactionTodayStats godoc
// @Summary 获取今日交易统计
// @Description 获取指定交易类型的今日统计数据
// @Tags 交易统计
// @Accept json
// @Produce json
// @Success 200 {object} protocol.Result{data=services.TodayStats}
// @Router /transactions/today-stats [post]
func (t *MerchantAdmin) GetTransactionTodayStats(c *gin.Context) {
	lang := middleware.GetLanguage(c)

	// 绑定请求参数
	var req TodayStatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// 从上下文获取商户信息
	mid := middleware.GetMidFromContext(c)
	// 调用服务层获取统计数据
	stats, code := services.GetMerchantTransactionService().GetTransactionTodayStats(mid, req.TrxType)
	if code != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(code, lang))
		return
	}

	// 返回成功结果
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(stats, lang))
}
