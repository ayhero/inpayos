package handlers

import (
	"inpayos/internal/middleware"
	"inpayos/internal/protocol"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (t *MerchantAdmin) Info(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	team := middleware.GetCashierTeamFromContext(c)
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(team.Protocol(), lang))
}
