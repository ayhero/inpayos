package handlers

import (
	"inpayos/internal/middleware"
	"inpayos/internal/protocol"
	"inpayos/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary 发送邮箱验证码
// @Description 发送邮箱验证码
// @Tags CashierAdmin
// @Accept json
// @Produce json
// @Param data body SendVerifyCodeReq true "发送验证码请求"
// @Success 200 {object} protocol.Result "返回结果"
// @Router /verifycode/send [post]
func SendVerifyCode(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	var req SendVerifyCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}
	// 发送验证码
	if err, _ := services.GetVerifyCodeService().SendEmailCode(strings.TrimSpace(req.Email), req.Type, lang); err != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(err, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(nil, lang))
}

// @Summary 验证邮箱验证码
// @Description 验证邮箱验证码是否正确
// @Tags CashierAdmin
// @Accept json
// @Produce json
// @Param request body VerifyCodeReq true "验证码验证请求"
// @Success 200 {object} protocol.Result "返回结果"
// @Router /verifycode/verify [post]
func VerifyCode(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	var req VerifyCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}
	if !services.GetVerifyCodeService().VerifyEmailCode(req.Type, strings.TrimSpace(req.Email), req.Code) {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(nil, lang))
}
