package handlers

import (
	"inpayos/internal/middleware"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/services"
	"inpayos/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Auth 授权认证
// @Summary      登陆授权认证
// @Description  处理登陆认证并返回token
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request  body      AuthRequest  true  "认证请求参数"
// @Router       /auth [post]
func (s *CashierAdmin) Auth(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	var merchant *models.Merchant

	if req.Email != "" {
		if req.Password == "" && req.Code == "" {
			c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.MissingParams, lang))
			return
		}
		if req.Password != "" && req.Code != "" {
			c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
			return
		}
		merchant = models.GetMerchantByEmail(req.Email)
		if merchant == nil {
			c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.MerchantNotFound, lang))
			return
		}

		merchant.Decrypt()
		if req.Code != "" {
			if merchant.GetG2FA() == "" {
				c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.TwoFactorRequired, lang))
				return
			}
			if req.Code == "" {
				c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.MissingParams, lang))
				return
			}
			if !services.VerifyG2FACode(merchant.GetG2FA(), req.Code) {
				c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidTwoFactorCode, lang))
				return
			}
		} else if req.Password != "" {
			if !merchant.IsPasswordValid(req.Password) {
				c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidCredentials, lang))
				return
			}
		}
	} else if req.Token != "" {
		jwtToken := middleware.ValidToken(c, []byte(req.Token))
		if jwtToken == nil || !jwtToken.Valid {
			c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidToken, lang))
			return
		}
		merchant = utils.GetContextData[models.Merchant](c, middleware.MerchantKey)
		if merchant == nil {
			c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.MerchantNotFound, lang))
			return
		}
	} else {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// Generate JWT token
	token, err := utils.GenerateMerchantTokenWithExpire(merchant.Mid, 0)
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.SystemError, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(AuthResponse{
		Token: token,
	}))
}

func (s *CashierAdmin) Logout(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(nil, lang))
}
