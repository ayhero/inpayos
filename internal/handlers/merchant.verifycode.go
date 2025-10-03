package handlers

import (
	"inpayos/internal/middleware"
	"inpayos/internal/protocol"
	"inpayos/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SendVerifyCodeReq 发送验证码请求
type SendVerifyCodeReq struct {
	Type  string `json:"type" binding:"required"`        // 验证码类型
	Email string `json:"email" binding:"required,email"` // 邮箱地址
}

// @Summary 发送邮箱验证码
// @Description 发送邮箱验证码
// @Tags 系统
// @Accept json
// @Produce json
// @Param data body VerifyCodeRequest true "邮箱"
// @Success 200 {object} protocol.Result "返回结果"
// @Router /verifycode/send [post]
func (t *MerchantAdmin) SendVerifyCode(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	var req SendVerifyCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}
	// 发送验证码
	if err := services.SendEmailVerifyCode(req.Type, strings.TrimSpace(req.Email)); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.SystemError, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(nil, lang))
}

// VerifyCodeReq 验证码验证请求
type VerifyCodeReq struct {
	Type  string `json:"type" binding:"required,oneof=register login reset"` // 验证码类型
	Email string `json:"email" binding:"required,email"`                     // 邮箱地址
	Code  string `json:"code" binding:"required"`                            // 验证码
}

// @Summary 验证邮箱验证码
// @Description 验证邮箱验证码是否正确
// @Tags 验证码
// @Accept json
// @Produce json
// @Param request body VerifyCodeReq true "验证码验证请求"
// @Success 200 {object} protocol.Result "返回结果"
// @Router /verifycode/verify [post]
func (t *MerchantAdmin) VerifyCode(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	var req VerifyCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}
	if !services.VerifyEmailCode(req.Type, req.Email, req.Code) {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(nil, lang))
}
