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
// @Summary 创建代收订单
// @Description 创建代收交易订单，支持多种支付方式和渠道
// @Tags OpenAPI
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.MerchantPayinRequest true "代收订单请求参数"
// @Success 200 {object} protocol.Result{data=protocol.Transaction} "创建成功"
// @Failure 400 {object} protocol.Result "请求参数错误"
// @Failure 401 {object} protocol.Result "认证失败"
// @Failure 500 {object} protocol.Result "服务器错误"
// @Router /openapi/payin [post]
func (a *OpenApi) Payin(c *gin.Context) {
	var req protocol.MerchantPayinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}
	lang := middleware.GetLanguage(c)
	req.Mid = middleware.GetMidFromContext(c)
	// 执行业务逻辑（代收类型）
	response, code := a.Transaction.CreatePayin(c, &req)
	result := protocol.HandleServiceResult(code, response, lang)
	c.JSON(http.StatusOK, result)
}

// Cancel 取消订单
// @Summary 取消订单
// @Description 取消指定的交易订单（代收或代付）
// @Tags OpenAPI
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.MerchantCancelRequest true "取消订单请求参数"
// @Success 200 {object} protocol.Result{data=protocol.Transaction} "取消成功"
// @Failure 400 {object} protocol.Result "请求参数错误"
// @Failure 401 {object} protocol.Result "认证失败"
// @Failure 500 {object} protocol.Result "服务器错误"
// @Router /openapi/cancel [post]
func (a *OpenApi) Cancel(c *gin.Context) {
	var req protocol.MerchantCancelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}
	lang := middleware.GetLanguage(c)
	req.Mid = middleware.GetMidFromContext(c)

	// 执行取消逻辑
	response, code := a.Transaction.Cancel(&req)
	result := protocol.HandleServiceResult(code, response, lang)
	c.JSON(http.StatusOK, result)
}

// Payout 创建代付订单
// @Summary 创建代付订单
// @Description 创建代付交易订单，向指定账户转账
// @Tags OpenAPI
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.MerchantPayoutRequest true "代付订单请求参数"
// @Success 200 {object} protocol.Result{data=protocol.Transaction} "创建成功"
// @Failure 400 {object} protocol.Result "请求参数错误"
// @Failure 401 {object} protocol.Result "认证失败"
// @Failure 500 {object} protocol.Result "服务器错误"
// @Router /openapi/payout [post]
func (a *OpenApi) Payout(c *gin.Context) {
	var req protocol.MerchantPayoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	lang := middleware.GetLanguage(c)
	req.Mid = middleware.GetMidFromContext(c)
	// 执行业务逻辑（代付类型）
	response, code := a.Transaction.CreatePayout(c, &req)
	result := protocol.HandleServiceResult(code, response, lang)
	c.JSON(http.StatusOK, result)
}

// =============================================================================
// 统一查询接口
// =============================================================================

// Query 查询交易状态/详情
// @Summary 查询交易状态
// @Description 根据请求ID或交易ID查询交易状态和详情
// @Tags OpenAPI
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.MerchantQueryRequest true "查询请求参数"
// @Success 200 {object} protocol.Result{data=protocol.Transaction} "查询成功"
// @Failure 400 {object} protocol.Result "请求参数错误"
// @Failure 401 {object} protocol.Result "认证失败"
// @Failure 500 {object} protocol.Result "服务器错误"
// @Router /openapi/query [post]
func (a *OpenApi) Query(c *gin.Context) {
	var req protocol.MerchantQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	lang := middleware.GetLanguage(c)
	req.Mid = middleware.GetMidFromContext(c)
	if req.ReqID == "" && req.TrxID == "" {
		lang := middleware.GetLanguage(c)
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.MissingParams, lang))
		return
	}
	response, code := a.Transaction.Query(&req)
	result := protocol.HandleServiceResult(code, response, lang)
	c.JSON(http.StatusOK, result)
}

// =============================================================================
// 余额查询
// =============================================================================

// Balance 查询账户余额
// @Summary 查询账户余额
// @Description 查询商户账户的各币种余额信息
// @Tags OpenAPI
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} protocol.Result{data=[]protocol.Account} "查询成功"
// @Failure 401 {object} protocol.Result "认证失败"
// @Failure 500 {object} protocol.Result "服务器错误"
// @Router /openapi/balance [post]
func (a *OpenApi) Balance(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	mid := middleware.GetMidFromContext(c)
	accounts := services.GetAccountService().GetMerchantAccountBalance(mid)
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(accounts, lang))
}
