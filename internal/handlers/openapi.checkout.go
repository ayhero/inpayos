package handlers

import (
	"inpayos/internal/middleware"
	"inpayos/internal/protocol"
	"net/http"

	"github.com/gin-gonic/gin"
)

// =============================================================================
// 收银台（Checkout）接口
// =============================================================================

// CreateCheckout 创建收银台会话
func (a *OpenApi) CreateCheckout(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	var req protocol.CreateCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// 从上下文获取收银台服务
	if a.Checkout == nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.SystemError, lang))
		return
	}

	// 执行业务逻辑
	response, code := a.Checkout.Create(&req)
	result := protocol.HandleServiceResult(code, response, lang)
	c.JSON(http.StatusOK, result)
}
func (a *OpenApi) ConfirmCheckout(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	var req protocol.ConfirmCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// 从上下文获取收银台服务
	if a.Checkout == nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.SystemError, lang))
		return
	}
	// 执行业务逻辑
	response, code := a.Checkout.Confirm(&req)
	c.JSON(http.StatusOK, protocol.HandleServiceResult(code, response, lang))
}

// GetCheckout 获取收银台信息
func (a *OpenApi) GetCheckout(c *gin.Context) {
	checkoutID := c.Query("checkout_id")
	lang := middleware.GetLanguage(c)
	if checkoutID == "" {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.MissingParams, lang))
		return
	}

	// 从上下文获取收银台服务
	if a.Checkout == nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.SystemError, lang))
		return
	}

	// 执行业务逻辑
	response, code := a.Checkout.Info(checkoutID)
	c.JSON(http.StatusOK, protocol.HandleServiceResult(code, response, lang))
}

// CancelCheckout 取消收银台会话
func (a *OpenApi) CancelCheckout(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	var req protocol.CancelCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}
	// 从上下文获取收银台服务
	if a.Checkout == nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.SystemError, lang))
		return
	}
	response, code := a.Checkout.Cancel(req.CheckoutID)
	c.JSON(http.StatusOK, protocol.HandleServiceResult(code, response, lang))
}
