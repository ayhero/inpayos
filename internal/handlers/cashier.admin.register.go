package handlers

import (
	"inpayos/internal/middleware"
	"inpayos/internal/protocol"
	"inpayos/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 商户注册
// @Description 邮箱注册新用户
// @Tags CashierAdmin
// @Accept json
// @Produce json
// @Param data body services.CashierTeamRegisterRequest true "注册信息"
// @Success 200 {object} protocol.Result "返回结果"
// @Router /register [post]
func (t *CashierAdmin) Register(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	var req services.CashierTeamRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}
	req.RegIP = c.ClientIP() // 获取注册IP
	// 调用服务层注册方法
	if err := services.RegisterCashierTeam(&req); err != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(err, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(nil, lang))
}
