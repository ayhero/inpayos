package handlers

import (
	"inpayos/internal/middleware"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListTransactions godoc
// @Summary 获取交易列表
// @Description 分页获取交易列表
// @Tags CashierAdmin
// @Accept json
// @Produce json
// @Success 200 {object} protocol.Result{data=protocol.PageResult{list=[]protocol.Transaction}}
// @Router /transactions/list [post]
func (t *CashierAdmin) ListTransactions(c *gin.Context) {
	lang := middleware.GetLanguage(c)

	// 绑定请求参数
	var req TransactionListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// 从上下文获取商户信息
	team := middleware.GetCashierTeamFromContext(c)

	// 构建查询参数
	query := &models.TrxQuery{
		Tid:              team.Tid,
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
	transactions, total, code := services.GetCashierTransactionService().ListTransactionByQuery(query)
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
// @Tags CashierAdmin
// @Accept json
// @Produce json
// @Success 200 {object} protocol.Result{data=protocol.Transaction}
// @Router /transactions/detail [post]
func (t *CashierAdmin) TransactionDetail(c *gin.Context) {
	lang := middleware.GetLanguage(c)

	// 绑定请求参数
	var req TransactionDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// 从上下文获取商户信息
	team := middleware.GetCashierTeamFromContext(c)

	// 构建查询参数
	query := &models.TrxQuery{
		Tid:     team.Tid,
		TrxType: req.TrxType,
		TrxID:   req.TrxID,
		Page:    1,
		Size:    1,
	}

	// 调用服务层查询交易详情
	transactions, _, code := services.GetCashierTransactionService().ListTransactionByQuery(query)
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

// GetTransactionTodayStats godoc
// @Summary 获取今日交易统计
// @Description 获取指定交易类型的今日统计数据
// @Tags CashierAdmin
// @Accept json
// @Produce json
// @Success 200 {object} protocol.Result{data=services.TodayStats}
// @Router /transactions/today-stats [post]
func (t *CashierAdmin) GetTransactionTodayStats(c *gin.Context) {
	lang := middleware.GetLanguage(c)

	// 绑定请求参数
	var req TodayStatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// 从上下文获取商户信息
	tid := middleware.GetTidFromContext(c)
	// 调用服务层获取统计数据
	stats, code := services.GetCashierTransactionService().GetTransactionTodayStats(tid, req.TrxType)
	if code != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(code, lang))
		return
	}

	// 返回成功结果
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(stats, lang))
}
