package services

import (
	"errors"
	"inpayos/internal/models"
	"inpayos/internal/protocol"
	"strconv"
	"strings"
)

type TrxRouterService struct {
}

// NewTrxRouterService 创建TrxRouter服务实例
func NewTrxRouterService() *TrxRouterService {
	return &TrxRouterService{}
}

// CreateTrxRouter 创建交易路由
func (s *TrxRouterService) CreateTrxRouter(req *protocol.TrxRouterRequest) (protocol.ErrorCode, *protocol.TrxRouterResponse) {
	// 验证三选一互斥条件
	if err := s.validateCashierSelection(req.CashierTypes, req.CashierIDs, req.CashierGroup); err != nil {
		return protocol.InvalidParams, nil
	}

	// 检查RouterID是否已存在
	if _, err := models.GetTrxRouterByID(req.RouterID); err == nil {
		return protocol.MerchantAlreadyExists, nil // 借用这个表示数据已存在
	}

	// 创建TrxRouter模型
	router := models.NewTrxRouter()
	router.RouterID = req.RouterID
	router.TrxRouterValues = s.buildTrxRouterValues(req)

	// 保存到数据库
	if err := models.CreateTrxRouter(router); err != nil {
		return protocol.SystemError, nil
	}

	return protocol.Success, s.convertToTrxRouterResponse(router)
}

// GetTrxRouter 获取交易路由详情
func (s *TrxRouterService) GetTrxRouter(routerID string) (protocol.ErrorCode, *protocol.TrxRouterResponse) {
	router, err := models.GetTrxRouterByID(routerID)
	if err != nil {
		return protocol.TransactionNotFound, nil
	}

	return protocol.Success, s.convertToTrxRouterResponse(router)
}

// UpdateTrxRouter 更新交易路由
func (s *TrxRouterService) UpdateTrxRouter(routerID string, req *protocol.TrxRouterUpdateRequest) (protocol.ErrorCode, *protocol.TrxRouterResponse) {
	// 获取现有路由
	router, err := models.GetTrxRouterByID(routerID)
	if err != nil {
		return protocol.TransactionNotFound, nil
	}

	// 验证三选一互斥条件
	cashierTypes := req.CashierTypes
	cashierIDs := req.CashierIDs
	cashierGroup := ""
	if req.CashierGroup != nil {
		cashierGroup = *req.CashierGroup
	}

	if err := s.validateCashierSelection(cashierTypes, cashierIDs, cashierGroup); err != nil {
		return protocol.InvalidParams, nil
	}

	// 更新字段
	s.updateTrxRouterFields(router, req)

	// 保存更新
	if err := models.UpdateTrxRouter(router); err != nil {
		return protocol.SystemError, nil
	}

	return protocol.Success, s.convertToTrxRouterResponse(router)
}

// DeleteTrxRouter 删除交易路由（软删除，设置状态为inactive）
func (s *TrxRouterService) DeleteTrxRouter(routerID string) protocol.ErrorCode {
	router, err := models.GetTrxRouterByID(routerID)
	if err != nil {
		return protocol.TransactionNotFound
	}

	status := "inactive"
	router.TrxRouterValues.Status = &status

	if err := models.UpdateTrxRouter(router); err != nil {
		return protocol.SystemError
	}

	return protocol.Success
}

// ListTrxRouters 获取交易路由列表
func (s *TrxRouterService) ListTrxRouters(req *protocol.TrxRouterListRequest) (protocol.ErrorCode, *protocol.TrxRouterListResponse) {
	// 构建查询条件
	query := models.DB.Model(&models.TrxRouter{})

	// 添加过滤条件
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.TransactionType != "" {
		query = query.Where("transaction_type = ?", req.TransactionType)
	}
	if req.PaymentMethod != "" {
		query = query.Where("payment_method = ?", req.PaymentMethod)
	}
	if req.Currency != "" {
		query = query.Where("currency = ?", req.Currency)
	}
	if req.Country != "" {
		query = query.Where("country = ?", req.Country)
	}
	if req.CashierType != "" {
		query = query.Where("JSON_CONTAINS(cashier_types, ?)", `"`+req.CashierType+`"`)
	}
	if req.CashierGroup != "" {
		query = query.Where("cashier_group = ?", req.CashierGroup)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return protocol.SystemError, nil
	}

	// 分页查询
	var routers []*models.TrxRouter
	offset := req.GetOffset()
	limit := req.GetLimit()

	if err := query.Order("priority DESC, created_at DESC").
		Offset(offset).Limit(limit).Find(&routers).Error; err != nil {
		return protocol.SystemError, nil
	}

	// 转换响应
	items := make([]*protocol.TrxRouterResponse, len(routers))
	for i, router := range routers {
		items[i] = s.convertToTrxRouterResponse(router)
	}

	pageResult := protocol.NewPageResult(items, total, &req.Pagination)
	return protocol.Success, &protocol.TrxRouterListResponse{PageResult: pageResult}
}

// MatchRoute 路由匹配
func (s *TrxRouterService) MatchRoute(req *protocol.RouteMatchRequest) (protocol.ErrorCode, *protocol.RouteMatchResponse) {
	// 如果指定了强制路由ID
	if req.ForceRouterID != "" {
		return s.matchByRouterID(req.ForceRouterID, req)
	}

	// 根据条件查找匹配的路由
	routers, err := models.GetTrxRoutersByCondition(
		req.TransactionType,
		req.PaymentMethod,
		req.Currency,
		req.Country,
		req.Amount,
	)
	if err != nil {
		return protocol.SystemError, nil
	}

	if len(routers) == 0 {
		return protocol.TransactionNotFound, nil
	}

	// 选择最佳路由
	bestRouter := s.selectBestRouter(routers, req)
	if bestRouter == nil {
		return protocol.TransactionNotFound, nil
	}

	// 获取可用的Cashier
	cashierInfo := s.getAvailableCashier(bestRouter, req)
	if cashierInfo == nil {
		return protocol.ChannelNotFound, nil
	}

	response := &protocol.RouteMatchResponse{
		Router:      s.convertToTrxRouterResponse(bestRouter),
		Cashier:     cashierInfo,
		MatchReason: "best_match",
		Metadata: map[string]interface{}{
			"total_candidates": len(routers),
			"selected_at":      models.GetCurrentTimeMillis(),
		},
	}

	return protocol.Success, response
}

// TestRoute 测试路由匹配
func (s *TrxRouterService) TestRoute(req *protocol.RouterTestRequest) (protocol.ErrorCode, *protocol.RouterTestResponse) {
	router, err := models.GetTrxRouterByID(req.RouterID)
	if err != nil {
		return protocol.TransactionNotFound, nil
	}

	// 测试匹配条件
	matched := router.TrxRouterValues.MatchesRequest(
		req.TransactionType,
		req.PaymentMethod,
		req.Currency,
		req.Country,
		req.Amount,
	)

	response := &protocol.RouterTestResponse{
		Matched: matched,
		MatchDetails: map[string]interface{}{
			"router_id":        req.RouterID,
			"transaction_type": req.TransactionType,
			"payment_method":   req.PaymentMethod,
			"amount":           req.Amount.String(),
			"currency":         req.Currency,
			"country":          req.Country,
		},
	}

	if !matched {
		response.FailedReasons = s.getMatchFailureReasons(router, req)
	} else {
		// 检查可用的Cashier
		matchReq := &protocol.RouteMatchRequest{
			TransactionType: req.TransactionType,
			PaymentMethod:   req.PaymentMethod,
			Amount:          req.Amount,
			Currency:        req.Currency,
			Country:         req.Country,
		}
		response.AvailableCashier = s.getAvailableCashier(router, matchReq)
	}

	return protocol.Success, response
}

// 辅助方法

// validateCashierSelection 验证Cashier选择的互斥性
func (s *TrxRouterService) validateCashierSelection(cashierTypes, cashierIDs []string, cashierGroup string) error {
	count := 0
	if len(cashierTypes) > 0 {
		count++
	}
	if len(cashierIDs) > 0 {
		count++
	}
	if cashierGroup != "" {
		count++
	}

	if count != 1 {
		return errors.New("CashierTypes, CashierIDs, and CashierGroup cannot have values simultaneously")
	}
	return nil
}

// buildTrxRouterValues 构建TrxRouterValues
func (s *TrxRouterService) buildTrxRouterValues(req *protocol.TrxRouterRequest) *models.TrxRouterValues {
	values := &models.TrxRouterValues{
		CashierTypes: models.StringArray{},
		CashierIDs:   models.StringArray{},
	}

	if req.Name != "" {
		values.SetName(req.Name)
	}
	if len(req.CashierTypes) > 0 {
		values.SetCashierTypes(req.CashierTypes)
	}
	if len(req.CashierIDs) > 0 {
		values.SetCashierIDs(req.CashierIDs)
	}
	if req.CashierGroup != "" {
		values.SetCashierGroup(req.CashierGroup)
	}
	if req.PaymentMethod != "" {
		values.SetPaymentMethod(req.PaymentMethod)
	}
	if req.TransactionType != "" {
		values.SetTransactionType(req.TransactionType)
	}
	if req.Currency != "" {
		values.SetCurrency(req.Currency)
	}
	if req.Country != "" {
		values.SetCountry(req.Country)
	}
	if req.MinAmount != nil {
		values.SetMinAmount(*req.MinAmount)
	}
	if req.MaxAmount != nil {
		values.SetMaxAmount(*req.MaxAmount)
	}
	if req.Priority > 0 {
		values.SetPriority(req.Priority)
	}
	if req.Status != "" {
		values.SetStatus(req.Status)
	}
	if req.EffectiveTime != nil {
		values.EffectiveTime = req.EffectiveTime
	}
	if req.ExpireTime != nil {
		values.ExpireTime = req.ExpireTime
	}
	if req.Remark != "" {
		values.SetRemark(req.Remark)
	}

	return values
}

// updateTrxRouterFields 更新TrxRouter字段
func (s *TrxRouterService) updateTrxRouterFields(router *models.TrxRouter, req *protocol.TrxRouterUpdateRequest) {
	if req.Name != nil {
		router.TrxRouterValues.SetName(*req.Name)
	}
	if len(req.CashierTypes) > 0 {
		router.TrxRouterValues.SetCashierTypes(req.CashierTypes)
	}
	if len(req.CashierIDs) > 0 {
		router.TrxRouterValues.SetCashierIDs(req.CashierIDs)
	}
	if req.CashierGroup != nil {
		router.TrxRouterValues.SetCashierGroup(*req.CashierGroup)
	}
	if req.PaymentMethod != nil {
		router.TrxRouterValues.SetPaymentMethod(*req.PaymentMethod)
	}
	if req.TransactionType != nil {
		router.TrxRouterValues.SetTransactionType(*req.TransactionType)
	}
	if req.Currency != nil {
		router.TrxRouterValues.SetCurrency(*req.Currency)
	}
	if req.Country != nil {
		router.TrxRouterValues.SetCountry(*req.Country)
	}
	if req.MinAmount != nil {
		router.TrxRouterValues.SetMinAmount(*req.MinAmount)
	}
	if req.MaxAmount != nil {
		router.TrxRouterValues.SetMaxAmount(*req.MaxAmount)
	}
	if req.Priority != nil {
		router.TrxRouterValues.SetPriority(*req.Priority)
	}
	if req.Status != nil {
		router.TrxRouterValues.SetStatus(*req.Status)
	}
	if req.EffectiveTime != nil {
		router.TrxRouterValues.EffectiveTime = req.EffectiveTime
	}
	if req.ExpireTime != nil {
		router.TrxRouterValues.ExpireTime = req.ExpireTime
	}
	if req.Remark != nil {
		router.TrxRouterValues.SetRemark(*req.Remark)
	}
}

// convertToTrxRouterResponse 转换为响应格式
func (s *TrxRouterService) convertToTrxRouterResponse(router *models.TrxRouter) *protocol.TrxRouterResponse {
	return &protocol.TrxRouterResponse{
		ID:              router.ID,
		RouterID:        router.RouterID,
		Name:            router.TrxRouterValues.GetName(),
		CashierTypes:    router.TrxRouterValues.GetCashierTypes(),
		CashierIDs:      router.TrxRouterValues.GetCashierIDs(),
		CashierGroup:    router.TrxRouterValues.GetCashierGroup(),
		PaymentMethod:   router.TrxRouterValues.GetPaymentMethod(),
		TransactionType: router.TrxRouterValues.GetTransactionType(),
		Currency:        router.TrxRouterValues.GetCurrency(),
		Country:         router.TrxRouterValues.GetCountry(),
		MinAmount:       router.TrxRouterValues.MinAmount,
		MaxAmount:       router.TrxRouterValues.MaxAmount,
		Priority:        router.TrxRouterValues.GetPriority(),
		Status:          router.TrxRouterValues.GetStatus(),
		EffectiveTime:   router.TrxRouterValues.EffectiveTime,
		ExpireTime:      router.TrxRouterValues.ExpireTime,
		Remark:          router.TrxRouterValues.GetRemark(),
		CreatedAt:       router.CreatedAt,
		UpdatedAt:       router.UpdatedAt,
	}
}

// matchByRouterID 根据路由ID强制匹配
func (s *TrxRouterService) matchByRouterID(routerID string, req *protocol.RouteMatchRequest) (protocol.ErrorCode, *protocol.RouteMatchResponse) {
	router, err := models.GetTrxRouterByID(routerID)
	if err != nil {
		return protocol.TransactionNotFound, nil
	}

	if !router.TrxRouterValues.IsEffective() {
		return protocol.ChannelDisabled, nil
	}

	cashierInfo := s.getAvailableCashier(router, req)
	if cashierInfo == nil {
		return protocol.ChannelNotFound, nil
	}

	response := &protocol.RouteMatchResponse{
		Router:      s.convertToTrxRouterResponse(router),
		Cashier:     cashierInfo,
		MatchReason: "forced_selection",
		Metadata: map[string]interface{}{
			"forced_router_id": routerID,
			"selected_at":      models.GetCurrentTimeMillis(),
		},
	}

	return protocol.Success, response
}

// selectBestRouter 选择最佳路由
func (s *TrxRouterService) selectBestRouter(routers []*models.TrxRouter, req *protocol.RouteMatchRequest) *models.TrxRouter {
	// 如果有偏好类型，优先选择匹配偏好的路由
	if len(req.PreferredTypes) > 0 {
		for _, router := range routers {
			for _, preferredType := range req.PreferredTypes {
				if router.TrxRouterValues.HasCashierType(preferredType) {
					return router
				}
			}
		}
	}

	// 选择优先级最高的路由
	if len(routers) > 0 {
		return routers[0] // 已按优先级排序
	}

	return nil
}

// getAvailableCashier 获取可用的Cashier
func (s *TrxRouterService) getAvailableCashier(router *models.TrxRouter, req *protocol.RouteMatchRequest) *protocol.CashierInfoResponse {
	selector := router.TrxRouterValues.GetCashierSelector()

	// 模拟Cashier选择逻辑
	switch selector["type"] {
	case "cashier_types":
		types := selector["values"].([]string)
		if len(types) > 0 {
			return &protocol.CashierInfoResponse{
				CashierID:   "CASH_" + strings.ToUpper(types[0]) + "_" + strconv.FormatInt(models.GetCurrentTimeMillis()%1000, 10),
				CashierType: types[0],
				Status:      "active",
				Available:   true,
				Config: map[string]interface{}{
					"selected_by": "type",
					"type":        types[0],
				},
			}
		}
	case "cashier_ids":
		ids := selector["values"].([]string)
		if len(ids) > 0 {
			return &protocol.CashierInfoResponse{
				CashierID:   ids[0],
				CashierType: "specific",
				Status:      "active",
				Available:   true,
				Config: map[string]interface{}{
					"selected_by": "id",
					"id":          ids[0],
				},
			}
		}
	case "cashier_group":
		group := selector["values"].(string)
		if group != "" {
			return &protocol.CashierInfoResponse{
				CashierID:   "CASH_" + strings.ToUpper(group) + "_" + strconv.FormatInt(models.GetCurrentTimeMillis()%1000, 10),
				CashierType: "group",
				Group:       group,
				Status:      "active",
				Available:   true,
				Config: map[string]interface{}{
					"selected_by": "group",
					"group":       group,
				},
			}
		}
	}

	return nil
}

// getMatchFailureReasons 获取匹配失败原因
func (s *TrxRouterService) getMatchFailureReasons(router *models.TrxRouter, req *protocol.RouterTestRequest) []string {
	var reasons []string

	// 检查各个匹配条件
	if router.TrxRouterValues.GetTransactionType() != "" &&
		router.TrxRouterValues.GetTransactionType() != "*" &&
		router.TrxRouterValues.GetTransactionType() != req.TransactionType {
		reasons = append(reasons, "Transaction type mismatch")
	}

	if router.TrxRouterValues.GetPaymentMethod() != "" &&
		router.TrxRouterValues.GetPaymentMethod() != "*" &&
		router.TrxRouterValues.GetPaymentMethod() != req.PaymentMethod {
		reasons = append(reasons, "Payment method mismatch")
	}

	if router.TrxRouterValues.GetCurrency() != "" &&
		router.TrxRouterValues.GetCurrency() != "*" &&
		router.TrxRouterValues.GetCurrency() != req.Currency {
		reasons = append(reasons, "Currency mismatch")
	}

	if router.TrxRouterValues.GetCountry() != "" &&
		router.TrxRouterValues.GetCountry() != "*" &&
		router.TrxRouterValues.GetCountry() != req.Country {
		reasons = append(reasons, "Country mismatch")
	}

	if !router.TrxRouterValues.IsInAmountRange(req.Amount) {
		reasons = append(reasons, "Amount out of range")
	}

	if !router.TrxRouterValues.IsEffective() {
		reasons = append(reasons, "Router not effective")
	}

	return reasons
}
