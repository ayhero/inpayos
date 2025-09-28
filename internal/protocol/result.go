package protocol

import "inpayos/internal/i18n"

// API响应结构 (保持向后兼容)
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Result 统一API响应结构 ()
type Result struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

// 常用的响应码
const (
	CODE_SUCCESS        = "0000" // 成功
	CODE_PARAM_ERROR    = "2001" // 参数错误
	CODE_AUTH_ERROR     = "3000" // 认证失败
	CODE_BUSINESS_ERROR = "1003" // 业务错误
	CODE_SYSTEM_ERROR   = "1000" // 系统错误
)

// NewSuccessResult 创建成功响应
func NewSuccessResult(data interface{}) *Result {
	return &Result{
		Code: CODE_SUCCESS,
		Msg:  "Success",
		Data: data,
	}
}

// NewErrorResult 创建错误响应
func NewErrorResult(code, message string) *Result {
	return &Result{
		Code: code,
		Msg:  message,
		Data: nil,
	}
}

// NewParamErrorResult 创建参数错误响应
func NewParamErrorResult(message string) *Result {
	return &Result{
		Code: CODE_PARAM_ERROR,
		Msg:  message,
		Data: nil,
	}
}

// NewBusinessErrorResult 创建业务错误响应
func NewBusinessErrorResult(message string) *Result {
	return &Result{
		Code: CODE_BUSINESS_ERROR,
		Msg:  message,
		Data: nil,
	}
}

// NewAuthErrorResult 创建认证错误响应
func NewAuthErrorResult() *Result {
	return &Result{
		Code: CODE_AUTH_ERROR,
		Msg:  "Authentication failed",
		Data: nil,
	}
}

// NewAuthErrorResultWithMsg 创建自定义消息的认证错误响应
func NewAuthErrorResultWithMsg(msg string) *Result {
	return &Result{
		Code: CODE_AUTH_ERROR,
		Msg:  msg,
		Data: nil,
	}
}

// NewSystemErrorResult 创建系统错误响应
func NewSystemErrorResult(message string) *Result {
	return &Result{
		Code: CODE_SYSTEM_ERROR,
		Msg:  message,
		Data: nil,
	}
}

// NewErrorResultWithCode 根据错误码创建错误响应
func NewErrorResultWithCode(code ErrorCode, lang string) *Result {
	message := i18n.Translate(code.GetCode(), lang)
	return &Result{
		Code: code.GetCode(),
		Msg:  message,
		Data: nil,
	}
}

// NewErrorResultWithCodeAndLang 根据错误码和语言创建错误响应
func NewErrorResultWithCodeAndLang(code ErrorCode, lang string, args ...interface{}) *Result {
	message := i18n.Translate(code.GetCode(), lang, args...)
	return &Result{
		Code: code.GetCode(),
		Msg:  message,
		Data: nil,
	}
}

// NewSuccessResultWithLang 创建多语言成功响应
func NewSuccessResultWithLang(data interface{}, lang string) *Result {
	message := i18n.Translate(Success.GetCode(), lang)
	return &Result{
		Code: Success.GetCode(),
		Msg:  message,
		Data: data,
	}
}

// HandleServiceResult 处理服务层响应（错误码和数据）
func HandleServiceResult(code ErrorCode, data interface{}, lang string) *Result {
	if code == Success {
		return NewSuccessResultWithLang(data, lang)
	}
	return NewErrorResultWithCode(code, lang)
}

// HandleServiceResultWithArgs 处理服务层响应（支持参数格式化）
func HandleServiceResultWithArgs(code ErrorCode, data interface{}, lang string, args ...interface{}) *Result {
	if code == Success {
		return NewSuccessResultWithLang(data, lang)
	}
	return NewErrorResultWithCodeAndLang(code, lang, args...)
}

// 分页请求参数 ()
type Pagination struct {
	Size int `json:"size" form:"size" binding:"min=1" example:"10"` // 每页记录数
	Page int `json:"page" form:"page" binding:"min=1" example:"1"`  // 当前页码，从1开始
}

// GetOffset 获取数据库查询的偏移量
func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.Size
}

// GetLimit 获取数据库查询的限制数
func (p *Pagination) GetLimit() int {
	return p.Size
}

// NewDefaultPagination 创建默认分页参数（第1页，每页10条）
func NewDefaultPagination() *Pagination {
	return &Pagination{
		Page: 1,
		Size: 10,
	}
}

// PageResult 分页数据结构 ()
type PageResult struct {
	ResultType string         `json:"result_type"` // 数据类型
	Size       int64          `json:"size"`        // 每页记录数
	Current    int64          `json:"current"`     // 当前页码
	Total      int64          `json:"total"`       // 总页数
	Count      int64          `json:"count"`       // 总记录数
	Records    any            `json:"records"`     // 记录列表
	Attach     map[string]any `json:"attach"`      // 附加数据
}

// CountPage 计算总页数
func (p *PageResult) CountPage() int64 {
	if p.Size == 0 {
		return 0
	}
	p.Count = p.Total / p.Size
	if p.Total%p.Size > 0 {
		p.Count++
	}
	return p.Count
}

// NewPageResult 创建分页数据结构
func NewPageResult(records any, total int64, page *Pagination) *PageResult {
	if page == nil {
		page = NewDefaultPagination()
	}
	count := total / int64(page.Size)
	if total%int64(page.Size) > 0 {
		count++
	}
	return &PageResult{
		Size:    int64(page.Size),
		Current: int64(page.Page),
		Total:   total,
		Count:   count,
		Records: records,
		Attach:  make(map[string]any),
	}
}

// NewSuccessPageResult 创建分页成功响应
func NewSuccessPageResult(records interface{}, total int64, page *Pagination) *Result {
	return NewSuccessResult(NewPageResult(records, total, page))
}

// AddAttach 添加附加数据
func (p *PageResult) AddAttach(key string, value interface{}) {
	if p.Attach == nil {
		p.Attach = make(map[string]interface{})
	}
	p.Attach[key] = value
}

// 保持向后兼容的分页结构
type PageRequest struct {
	Page     int `json:"page" form:"page" binding:"min=1"`
	PageSize int `json:"page_size" form:"page_size" binding:"min=1,max=100"`
}

// 保持向后兼容的分页响应
type PageResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}
