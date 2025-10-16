package handlers

import (
	"inpayos/internal/middleware"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"inpayos/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 这个文件扩展了CashierAdmin结构体，添加了出纳员管理相关的方法
// CashierAdmin结构体定义在cashier.admin.go中

// CashierListRequest 出纳员列表请求
type CashierListRequest struct {
	Tid            string `json:"tid"`                          // 团队ID
	Name           string `json:"name"`                         // 出纳员姓名
	Email          string `json:"email"`                        // 邮箱
	Phone          string `json:"phone"`                        // 电话
	Status         string `json:"status"`                       // 状态
	Type           string `json:"type"`                         // 类型
	Region         string `json:"region"`                       // 地区
	CreatedAtStart int64  `json:"created_at_start"`             // 创建开始时间
	CreatedAtEnd   int64  `json:"created_at_end"`               // 创建结束时间
	Page           int    `json:"page" binding:"min=1"`         // 页码
	Size           int    `json:"size" binding:"min=1,max=100"` // 每页记录数
}

// CashierDetailRequest 出纳员详情请求
type CashierDetailRequest struct {
	Tid string `json:"tid" binding:"required"` // 团队ID
}

// ListCashiers godoc
// @Summary 获取出纳员列表
// @Description 分页获取出纳员列表
// @Tags CashierAdmin
// @Accept json
// @Produce json
// @Success 200 {object} protocol.Result{data=protocol.PageResult{list=[]protocol.CashierTeam}}
// @Router /cashiers/list [post]
func (t *CashierAdmin) ListCashiers(c *gin.Context) {
	lang := middleware.GetLanguage(c)

	// 绑定请求参数
	var req CashierListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// 从上下文获取当前出纳员团队信息
	team := middleware.GetCashierTeamFromContext(c)

	// 构建查询参数
	query := &models.CashierTeamQuery{
		Tid:            team.Tid, // 当前团队的Tid
		Name:           req.Name,
		Email:          req.Email,
		Phone:          req.Phone,
		Status:         req.Status,
		Type:           req.Type,
		Region:         req.Region,
		CreatedAtStart: req.CreatedAtStart,
		CreatedAtEnd:   req.CreatedAtEnd,
		Page:           req.Page,
		Size:           req.Size,
	}

	// 调用服务层查询出纳员列表
	cashiers, total, code := services.GetCashierAdminService().ListCashiersByQuery(query)
	if code != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(code, lang))
		return
	}

	// 转换为协议格式
	var list []*protocol.CashierTeam
	for _, cashier := range cashiers {
		list = append(list, cashier.Protocol())
	}

	// 构建分页结果
	pagination := &protocol.Pagination{
		Page: req.Page,
		Size: req.Size,
	}

	// 返回成功结果
	c.JSON(http.StatusOK, protocol.NewSuccessPageResult(list, total, pagination))
}

// CashierDetail godoc
// @Summary 获取出纳员详情
// @Description 获取指定出纳员的详细信息
// @Tags CashierAdmin
// @Accept json
// @Produce json
// @Success 200 {object} protocol.Result{data=protocol.CashierTeam}
// @Router /cashiers/detail [post]
func (t *CashierAdmin) CashierDetail(c *gin.Context) {
	lang := middleware.GetLanguage(c)

	// 绑定请求参数
	var req CashierDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.InvalidParams, lang))
		return
	}

	// 从上下文获取当前出纳员团队信息
	team := middleware.GetCashierTeamFromContext(c)

	// 如果请求的Tid与当前团队Tid不同，则检查权限
	if req.Tid != team.Tid {
		// 这里可以添加更多的权限检查逻辑
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(protocol.PermissionDenied, lang))
		return
	}

	// 调用服务层查询出纳员详情
	cashier, code := services.GetCashierAdminService().GetCashierDetail(req.Tid)
	if code != protocol.Success {
		c.JSON(http.StatusOK, protocol.NewErrorResultWithCode(code, lang))
		return
	}

	// 转换为协议格式并返回
	cashierInfo := cashier.Protocol()
	c.JSON(http.StatusOK, protocol.NewSuccessResultWithLang(cashierInfo, lang))
}
