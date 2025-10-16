package handlers

import (
	"inpayos/internal/middleware"
	"inpayos/internal/protocol"
	"inpayos/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTodayStats godoc
// @Summary 获取今日统计数据
// @Description 获取今日代收、代付、成功率、结算、余额等关键指标
// @Tags CashierAdmin
// @Accept json
// @Produce json
// @Param request body DashboardTodayStatsRequest true "请求参数"
// @Success 200 {object} protocol.Result{data=services.DashboardTodayStats}
// @Router /dashboard/today-stats [post]
func (m *CashierAdmin) GetTodayStats(c *gin.Context) {
	var req DashboardTodayStatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult(err.Error()))
		return
	}

	tid := middleware.GetTidFromContext(c)
	// 调用服务层获取今日统计数据
	stats, err := services.GetCashierTeamTodayStats(tid, req.Currency)
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewBusinessErrorResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(stats))
}

// GetDashboardOverview godoc
// @Summary 获取Dashboard概览数据
// @Description 一次性获取Dashboard所有模块的数据，减少前端请求次数
// @Tags CashierAdmin
// @Accept json
// @Produce json
// @Success 200 {object} protocol.Result{data=services.DashboardOverview}
// @Router /dashboard/overview [post]
func (m *CashierAdmin) GetDashboardOverview(c *gin.Context) {
	// 获取当前商户信息
	tid := middleware.GetTidFromContext(c)

	// 调用服务层获取概览数据
	overview, err := services.GetCashierTeamDashboardOverview(tid)
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewBusinessErrorResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(overview))
}

// GetAccountBalance godoc
// @Summary 获取账户余额
// @Description 获取商户各币种账户余额信息
// @Tags CashierAdmin
// @Accept json
// @Produce json
// @Success 200 {object} protocol.Result{data=[]services.DashboardAccountBalance}
// @Router /dashboard/account-balance [post]
func (m *CashierAdmin) GetAccountBalance(c *gin.Context) {
	// 获取当前商户信息
	merchant := middleware.GetMerchantFromContext(c)
	if merchant == nil {
		c.JSON(http.StatusOK, protocol.NewAuthErrorResult())
		return
	}
	merchantID := merchant.ID

	// 调用服务层获取账户余额
	balances, err := services.GetCashierTeamAccountBalance(merchantID)
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewBusinessErrorResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(balances))
}
