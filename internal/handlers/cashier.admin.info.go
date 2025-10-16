package handlers

import (
	"inpayos/internal/middleware"
	"inpayos/internal/protocol"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (t *CashierAdmin) Info(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	user := middleware.GetMerchantFromContext(c)
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(user.Protocol(), lang))
}
