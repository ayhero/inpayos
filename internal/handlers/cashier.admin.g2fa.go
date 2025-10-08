package handlers

import (
	"fmt"
	"inpayos/internal/middleware"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 商户绑定G2FA
// @Description 绑定商户的二次验证
// @Tags 商户API
// @Accept json
// @Produce json
// @Param data body BindG2FAReq true "绑定信息"
// @Success 200 {object} protocol.Result "返回结果"
// @Router /merchant/g2fa/bind [post]
func (t *CashierAdmin) BindG2FA(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	var req BindG2FAReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// 获取当前商户信息
	cashierTeam := middleware.GetCashierTeamFromContext(c)
	//已经绑定过，重新绑定
	if cashierTeam.GetG2FA() != "" {
		// 验证码校验
		if !services.VerifyEmailCode(protocol.VerifyCodeTypeResetG2FA, cashierTeam.GetEmail(), req.VerifyCode) {
			c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
			return
		}
	}

	// 从缓存获取待绑定的G2FA密钥
	cacheKey := fmt.Sprintf(protocol.G2FABindingTpl, cashierTeam.Tid)
	newG2FAKey, err := models.GetCache(cacheKey)
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.SystemError, lang))
		return
	}

	// 使用待绑定的新密钥验证G2FA code
	if !services.VerifyG2FACode(newG2FAKey, req.Code) {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidTwoFactorCode, lang))
		return
	}

	// 验证通过后，更新商户的G2FA信息
	if err := models.WriteDB.Model(cashierTeam).Updates(&models.MerchantValues{G2FA: &newG2FAKey}).Error; err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.DatabaseError, lang))
		return
	}

	// 删除缓存中的临时G2FA密钥
	models.Delete(cacheKey)

	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(nil, lang))
}

// @Summary 生成新的G2FA密钥
// @Description 为商户生成新的G2FA密钥
// @Tags 商户API
// @Accept json
// @Produce json
// @Success 200 {object} protocol.Result "返回结果"
// @Router /merchant/g2fa/new [post]
func (t *CashierAdmin) NewG2FA(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	// 获取当前商户信息
	merchant := middleware.GetMerchantFromContext(c)
	// 生成新的G2FA密钥
	newG2FAKey := services.GenerateG2FAKey()
	if newG2FAKey == "" {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.SystemError, lang))
		return
	}

	// 将新生成的G2FA密钥存入缓存
	cacheKey := fmt.Sprintf(protocol.G2FABindingTpl, merchant.Mid)
	if err := models.SetCache(cacheKey, newG2FAKey, protocol.G2FACacheExpiration); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.CacheError, lang))
		return
	}

	response := G2FAResponse{
		G2FAKey: newG2FAKey,
		QRCode:  services.GenerateG2FAQRCode(merchant.Mid, newG2FAKey),
	}
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(response, lang))
}
