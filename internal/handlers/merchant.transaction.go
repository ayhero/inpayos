package handlers

import (
	"github.com/gin-gonic/gin"
)

// ListTransactions godoc
// @Summary 获取交易列表
// @Description 分页获取交易列表
// @Tags 交易管理
// @Accept json
// @Produce json
// @Param page query int true "页码"
// @Param size query int true "每页数量"
// @Param start_time query int false "开始时间"
// @Param end_time query int false "结束时间"
// @Param trx_id query string false "交易ID"
// @Param req_id query string false "商户订单号"
// @Param trx_type query string true "交易类型"
// @Success 200 {object} protocol.Result{data=protocol.PageResult{list=[]protocol.TransactionInfo}}
// @Router /transactions [get]
func (t *MerchantAdmin) ListTransactions(c *gin.Context) {

}

// GetTransaction godoc
// @Summary 获取交易详情
// @Description 获取指定交易ID的详细信息
// @Tags 交易管理
// @Accept json
// @Produce json
// @Param trx_id query string true "交易ID"
// @Success 200 {object} protocol.Result{data=protocol.TransactionInfo}
// @Router /transactions/detail [get]
func (t *MerchantAdmin) GetTransaction(c *gin.Context) {
}
