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
// @Summary 创建收银台会话
// @Description 创建一个新的收银台会话，用于集成支付页面
// @Tags OpenAPI
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body protocol.CreateCheckoutRequest true "创建收银台请求参数"
// @Success 200 {object} protocol.Result{data=protocol.Checkout} "创建成功"
// @Failure 400 {object} protocol.Result "请求参数错误"
// @Failure 401 {object} protocol.Result "认证失败"
// @Failure 500 {object} protocol.Result "服务器错误"
// @Router /checkout [post]
func (a *OpenApi) CreateCheckout(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	var req protocol.CreateCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}
	// 获取商户ID
	req.Mid = middleware.GetMidFromContext(c)
	// 执行业务逻辑
	response, code := a.Checkout.Create(&req)
	c.JSON(http.StatusOK, protocol.HandleServiceResult(code, response, lang))
}
