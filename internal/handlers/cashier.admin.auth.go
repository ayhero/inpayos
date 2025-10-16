package handlers

import (
	"inpayos/internal/config"
	"inpayos/internal/middleware"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Auth 授权认证
// @Summary      登陆授权认证
// @Description  处理登陆认证并返回token
// @Tags         CashierAdmin
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

	var cashier *models.CashierTeam

	if req.Email != "" {
		if req.Password == "" && req.Code == "" {
			c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.MissingParams, lang))
			return
		}
		if req.Password != "" && req.Code != "" {
			c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
			return
		}
		cashier = models.GetCashierTeamByEmail(req.Email)
		if cashier == nil {
			c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.MerchantNotFound, lang))
			return
		}

		cashier.Decrypt()
		if req.Code != "" {
			if cashier.GetG2FA() == "" {
				c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.TwoFactorRequired, lang))
				return
			}
			if req.Code == "" {
				c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.MissingParams, lang))
				return
			}
			if !services.VerifyG2FACode(cashier.GetG2FA(), req.Code) {
				c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidTwoFactorCode, lang))
				return
			}
		} else if req.Password != "" {
			if !cashier.IsPasswordValid(req.Password) {
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
		cashier = middleware.GetCashierTeamFromContext(c)
		if cashier == nil {
			c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.MerchantNotFound, lang))
			return
		}
	} else {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(cashier.Tid, time.Now().Add(72*time.Hour), config.Get().Server.CashierAdmin.Jwt.Secret)
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.SystemError, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(AuthResponse{
		Token: token,
	}))
}

// @Summary 收银员登出
// @Description 收银员登出接口
// @Tags CashierAdmin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} protocol.Result "登出成功"
// @Router /logout [post]
func (s *CashierAdmin) Logout(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(nil, lang))
}
