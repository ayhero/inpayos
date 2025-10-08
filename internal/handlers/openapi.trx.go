package handlers

import (
	"inpayos/internal/middleware"
	"inpayos/internal/protocol"
	"inpayos/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// =============================================================================
// 代收（Payin）接口
// =============================================================================

// Payin 创建代收订单
func (a *OpenApi) Payin(c *gin.Context) {
	var req protocol.MerchantPayinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// 执行业务逻辑（代收类型）
	response, code := a.Transaction.CreatePayin(&req)
	lang := middleware.GetLanguage(c)
	result := protocol.HandleServiceResult(code, response, lang)
	c.JSON(http.StatusOK, result)
}

// Cancel 取消订单
func (a *OpenApi) Cancel(c *gin.Context) {
	var req protocol.MerchantCancelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// 执行取消逻辑
	response, code := a.Transaction.Cancel(&req)
	lang := middleware.GetLanguage(c)
	result := protocol.HandleServiceResult(code, response, lang)
	c.JSON(http.StatusOK, result)
}

// Payout 创建代付订单
func (a *OpenApi) Payout(c *gin.Context) {
	var req protocol.MerchantPayoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// 执行业务逻辑（代付类型）
	response, code := a.Transaction.CreatePayout(&req)
	lang := middleware.GetLanguage(c)
	result := protocol.HandleServiceResult(code, response, lang)
	c.JSON(http.StatusOK, result)
}

// =============================================================================
// 统一查询接口
// =============================================================================

// Query 查询交易状态/详情
func (a *OpenApi) Query(c *gin.Context) {
	var req protocol.MerchantQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	if req.ReqID == "" && req.TrxID == "" {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.MissingParams, lang))
		return
	}
	response, code := a.Transaction.Query(&req)
	lang := middleware.GetLanguage(c)
	result := protocol.HandleServiceResult(code, response, lang)
	c.JSON(http.StatusOK, result)
}

// =============================================================================
// 余额查询
// =============================================================================

// Balance 查询账户余额
func (a *OpenApi) Balance(c *gin.Context) {
	accountService := services.GetAccountService()
	mid := middleware.GetMidFromContext(c)
	balance, code := accountService.GetMerchantAccountBalance(mid)
	lang := middleware.GetLanguage(c)
	if code != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(code, lang))
		return
	}
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(balance, lang))
}
