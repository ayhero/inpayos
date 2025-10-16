package handlers

import (
	"inpayos/internal/middleware"
	"inpayos/internal/protocol"
	"inpayos/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ResetPasswordReq 重置密码请求
type ResetPasswordReq struct {
	Email            string `json:"email" binding:"required,email"`       // 邮箱地址
	VerificationCode string `json:"verification_code" binding:"required"` // 验证码
}

// @Summary 重置密码
// @Description 通过邮箱验证码重置密码，新密码将发送到邮箱
// @Tags Merchant
// @Accept json
// @Produce json
// @Param data body ResetPasswordReq true "重置密码请求"
// @Success 200 {object} protocol.Result "返回结果"
// @Router /password/reset [post]
func (t *MerchantAdmin) ResetPassword(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	var req ResetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	email := strings.TrimSpace(req.Email)
	code := strings.TrimSpace(req.VerificationCode)

	// 验证验证码
	if !services.GetVerifyCodeService().VerifyEmailCode(protocol.VerifyCodeTypeResetPassword, email, code) {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// 重置密码
	newPassword, err := services.ResetMerchantPassword(email)
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.SystemError, lang))
		return
	}

	// 发送新密码到邮箱
	if err := services.SendNewPasswordEmail(email, newPassword); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.SystemError, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(nil, lang))
}

// ChangePasswordReq 修改密码请求
type ChangePasswordReq struct {
	NewPassword string `json:"new_password" binding:"required,min=6"` // 新密码
}

// @Summary 修改密码
// @Description 商户登录后修改密码
// @Tags Merchant
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param data body ChangePasswordReq true "修改密码请求"
// @Success 200 {object} protocol.Result "返回结果"
// @Router /password/change [post]
func (t *MerchantAdmin) ChangePassword(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	var req ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// 从上下文中获取商户信息
	merchant := middleware.GetMerchantFromContext(c)
	// 修改密码
	if err := services.ChangeMerchantPassword(merchant.GetEmail(), req.NewPassword); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.SystemError, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(nil, lang))
}
