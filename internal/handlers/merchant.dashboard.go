package handlers

import (
	"inpayos/internal/middleware"
	"inpayos/internal/protocol"
	"inpayos/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DashboardTodayStatsRequest 今日统计请求
type DashboardTodayStatsRequest struct {
	// 可以添加筛选条件，比如币种等
	Currency string `json:"currency" form:"currency"` // 币种筛选，可选
}

// DashboardTransactionTrendRequest 交易趋势请求
type DashboardTransactionTrendRequest struct {
	Days int `json:"days" form:"days" binding:"min=1,max=30"` // 查询天数，默认7天
}

// DashboardSettlementTrendRequest 结算趋势请求
type DashboardSettlementTrendRequest struct {
	Weeks int `json:"weeks" form:"weeks" binding:"min=1,max=12"` // 查询周数，默认4周
}

// GetTodayStats godoc
// @Summary 获取今日统计数据
// @Description 获取今日代收、代付、成功率、结算、余额等关键指标
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param request body DashboardTodayStatsRequest true "请求参数"
// @Success 200 {object} protocol.Result{data=services.DashboardTodayStats}
// @Router /dashboard/today-stats [post]
func (m *MerchantAdmin) GetTodayStats(c *gin.Context) {
	var req DashboardTodayStatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult(err.Error()))
		return
	}

	// 获取当前商户信息
	merchant := middleware.GetMerchantFromContext(c)
	if merchant == nil {
		c.JSON(http.StatusOK, protocol.NewAuthErrorResult())
		return
	}
	merchantID := merchant.ID

	// 调用服务层获取今日统计数据
	stats, err := services.GetTodayStats(merchantID, req.Currency)
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewBusinessErrorResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(stats))
}

// GetTransactionTrend godoc
// @Summary 获取交易趋势数据
// @Description 获取指定天数内的代收、代付交易趋势数据
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param request body DashboardTransactionTrendRequest true "请求参数"
// @Success 200 {object} protocol.Result{data=[]services.DashboardTransactionTrend}
// @Router /dashboard/transaction-trend [post]
func (m *MerchantAdmin) GetTransactionTrend(c *gin.Context) {
	var req DashboardTransactionTrendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult(err.Error()))
		return
	}

	// 设置默认值
	if req.Days == 0 {
		req.Days = 7
	}

	// 获取当前商户信息
	merchant := middleware.GetMerchantFromContext(c)
	if merchant == nil {
		c.JSON(http.StatusOK, protocol.NewAuthErrorResult())
		return
	}
	merchantID := merchant.ID

	// 调用服务层获取交易趋势数据
	trends, err := services.GetTransactionTrend(merchantID, req.Days)
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewBusinessErrorResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(trends))
}

// GetSettlementTrend godoc
// @Summary 获取结算趋势数据
// @Description 获取指定周数内的结算趋势数据
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param request body DashboardSettlementTrendRequest true "请求参数"
// @Success 200 {object} protocol.Result{data=[]services.DashboardSettlementTrend}
// @Router /dashboard/settlement-trend [post]
func (m *MerchantAdmin) GetSettlementTrend(c *gin.Context) {
	var req DashboardSettlementTrendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewParamErrorResult(err.Error()))
		return
	}

	// 设置默认值
	if req.Weeks == 0 {
		req.Weeks = 4
	}

	// 获取当前商户信息
	merchant := middleware.GetMerchantFromContext(c)
	if merchant == nil {
		c.JSON(http.StatusOK, protocol.NewAuthErrorResult())
		return
	}
	merchantID := merchant.ID

	// 调用服务层获取结算趋势数据
	trends, err := services.GetSettlementTrend(merchantID, req.Weeks)
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewBusinessErrorResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(trends))
}

// GetDashboardOverview godoc
// @Summary 获取Dashboard概览数据
// @Description 一次性获取Dashboard所有模块的数据，减少前端请求次数
// @Tags Dashboard
// @Accept json
// @Produce json
// @Success 200 {object} protocol.Result{data=services.DashboardOverview}
// @Router /dashboard/overview [post]
func (m *MerchantAdmin) GetDashboardOverview(c *gin.Context) {
	// 获取当前商户信息
	merchant := middleware.GetMerchantFromContext(c)
	if merchant == nil {
		c.JSON(http.StatusOK, protocol.NewAuthErrorResult())
		return
	}
	merchantID := merchant.ID

	// 调用服务层获取概览数据
	overview, err := services.GetDashboardOverview(merchantID)
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewBusinessErrorResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(overview))
}

// GetAccountBalance godoc
// @Summary 获取账户余额
// @Description 获取商户各币种账户余额信息
// @Tags Dashboard
// @Accept json
// @Produce json
// @Success 200 {object} protocol.Result{data=[]services.DashboardAccountBalance}
// @Router /dashboard/account-balance [post]
func (m *MerchantAdmin) GetAccountBalance(c *gin.Context) {
	// 获取当前商户信息
	merchant := middleware.GetMerchantFromContext(c)
	if merchant == nil {
		c.JSON(http.StatusOK, protocol.NewAuthErrorResult())
		return
	}
	merchantID := merchant.ID

	// 调用服务层获取账户余额
	balances, err := services.GetAccountBalance(merchantID)
	if err != nil {
		c.JSON(http.StatusOK, protocol.NewBusinessErrorResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, protocol.NewSuccessResult(balances))
}
