package handlers

import (
	"inpayos/internal/middleware"
	"inpayos/internal/protocol"
	"inpayos/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AccountList 获取账户列表
func (t *CashierAdmin) AccountList(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	tid := middleware.GetTidFromContext(c)
	// 调用服务层
	accountService := services.GetAccountService()
	accounts, code := accountService.GetAccountList(tid, protocol.UserTypeCashierTeam)
	if code != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(code, lang))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(accounts, lang))
}

// AccountFlowList 获取账户流水列表
func (t *CashierAdmin) AccountFlowList(c *gin.Context) {
	lang := middleware.GetLanguage(c)
	tid := middleware.GetTidFromContext(c)

	// 绑定请求参数
	var req protocol.AccountFlowListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// 调用服务层查询流水列表
	accountService := services.GetAccountService()
	flows, total, code := accountService.ListAccountFlowByQuery(tid, protocol.UserTypeCashierTeam, &req)
	if code != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(code, lang))
		return
	}

	// 转换为协议格式
	var list []*protocol.FundFlow
	for _, flow := range flows {
		list = append(list, flow.Protocol())
	}

	// 构建分页结果
	pagination := &protocol.Pagination{
		Page: req.Page,
		Size: req.Size,
	}

	// 返回成功结果
	c.JSON(http.StatusOK, protocol.NewSuccessPageResult(list, total, pagination))
}
