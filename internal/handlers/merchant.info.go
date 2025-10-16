package handlers

import (
	"inpayos/internal/middleware"
	"inpayos/internal/protocol"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (t *MerchantAdmin) Info(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	merchant := middleware.GetMerchantFromContext(c)
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(merchant.Protocol(), lang))
}
