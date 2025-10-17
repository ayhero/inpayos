package handlers

import (
	"inpayos/internal/middleware"
	"inpayos/internal/protocol"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SubmitCheckout 提交收银台支付
// @Summary 提交收银台支付
// @Description 用户在收银台页面提交支付信息，进行支付处理
// @Tags OpenAPI
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.SubmitCheckoutRequest true "提交支付请求参数"
// @Success 200 {object} protocol.Result{data=protocol.Checkout} "提交成功"
// @Failure 400 {object} protocol.Result "请求参数错误"
// @Failure 401 {object} protocol.Result "认证失败"
// @Failure 500 {object} protocol.Result "服务器错误"
// @Router /checkout/submit [post]
func (a *MerchantAdmin) SubmitCheckout(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	var req protocol.SubmitCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}
	// 获取商户ID
	req.Mid = middleware.GetMidFromContext(c)
	// 执行业务逻辑
	response, code := a.Checkout.Submit(&req)
	c.JSON(http.StatusOK, protocol.HandleServiceResult(code, response, lang))
}

// CheckoutServices 获取收银台可用的服务列表
// @Summary 获取收银台可用的服务列表
// @Description 根据收银台ID获取该收银台会话可用的支付服务和通道
// @Tags OpenAPI
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param checkout_id query string true "收银台会话ID"
// @Success 200 {object} protocol.Result{data=[]protocol.CheckoutService} "获取成功"
// @Failure 400 {object} protocol.Result "请求参数错误"
// @Failure 401 {object} protocol.Result "认证失败"
// @Failure 500 {object} protocol.Result "服务器错误"
// @Router /checkout/services [get]
func (a *MerchantAdmin) CheckoutServices(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	var req protocol.CheckoutInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}
	// 执行业务逻辑
	response, code := a.Checkout.Configs(req.CheckoutID)
	c.JSON(http.StatusOK, protocol.HandleServiceResult(code, response, lang))
}

// ConfirmCheckout 确认收银台支付
// @Summary 确认收银台支付
// @Description 确认用户在收银台的支付操作
// @Tags OpenAPI
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.ConfirmCheckoutRequest true "确认支付请求参数"
// @Success 200 {object} protocol.Result{data=protocol.Checkout} "确认成功"
// @Failure 400 {object} protocol.Result "请求参数错误"
// @Failure 401 {object} protocol.Result "认证失败"
// @Failure 500 {object} protocol.Result "服务器错误"
// @Router /checkout/confirm [post]
func (a *MerchantAdmin) ConfirmCheckout(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	var req protocol.ConfirmCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}
	req.Mid = middleware.GetMidFromContext(c)
	// 执行业务逻辑
	response, code := a.Checkout.Confirm(&req)
	c.JSON(http.StatusOK, protocol.HandleServiceResult(code, response, lang))
}

// CheckoutInfo 获取收银台信息
// @Summary 获取收银台信息
// @Description 根据收银台ID获取收银台会话的详细信息
// @Tags OpenAPI
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param checkout_id query string true "收银台会话ID"
// @Success 200 {object} protocol.Result{data=protocol.Checkout} "获取成功"
// @Failure 400 {object} protocol.Result "请求参数错误"
// @Failure 401 {object} protocol.Result "认证失败"
// @Failure 500 {object} protocol.Result "服务器错误"
// @Router /checkout/info [post]
func (a *MerchantAdmin) CheckoutInfo(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	var req protocol.CheckoutInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}
	// 执行业务逻辑
	response, code := a.Checkout.Info(req.CheckoutID)
	c.JSON(http.StatusOK, protocol.HandleServiceResult(code, response, lang))
}

// CancelCheckout 取消收银台会话
// @Summary 取消收银台会话
// @Description 取消指定的收银台会话
// @Tags OpenAPI
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.CancelCheckoutRequest true "取消收银台请求参数"
// @Success 200 {object} protocol.Result{data=protocol.Checkout} "取消成功"
// @Failure 400 {object} protocol.Result "请求参数错误"
// @Failure 401 {object} protocol.Result "认证失败"
// @Failure 500 {object} protocol.Result "服务器错误"
// @Router /checkout/cancel [post]
func (a *MerchantAdmin) CancelCheckout(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	var req protocol.CancelCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}
	req.Mid = middleware.GetMidFromContext(c)
	response, code := a.Checkout.Cancel(req.CheckoutID)
	c.JSON(http.StatusOK, protocol.HandleServiceResult(code, response, lang))
}
