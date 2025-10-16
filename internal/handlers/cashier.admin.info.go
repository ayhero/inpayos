package handlers

import (
	"inpayos/internal/middleware"
	"inpayos/internal/protocol"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 获取收银员信息
// @Description 获取当前登录收银员的详细信息
// @Tags CashierAdmin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} protocol.Result{data=protocol.CashierTeam} "获取成功"
// @Failure 401 {object} protocol.Result "认证失败"
// @Failure 500 {object} protocol.Result "服务器错误"
// @Router /info [post]
func (t *CashierAdmin) Info(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	user := middleware.GetCashierTeamFromContext(c)
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(user.Protocol(), lang))
}
