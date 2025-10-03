package protocol

// Page 分页请求参数和响应信息
type Page struct {
	Size      int    `json:"size" form:"size" binding:"min=1" example:"10"` // 每页记录数
	Page      int    `json:"page" form:"page" binding:"min=1" example:"1"`  // 当前页码，从1开始
	HasMore   bool   `json:"has_more,omitempty"`                            // 是否有更多数据（响应时使用）
	RequestID string `json:"request_id,omitempty"`                          // 请求ID（响应时使用）
}

// GetOffset 获取数据库查询的偏移量
func (p *Page) GetOffset() int {
	return (p.Page - 1) * p.Size
}

// GetLimit 获取数据库查询的限制数
func (p *Page) GetLimit() int {
	return p.Size
}

// NewDefaultPage 创建默认分页参数（第1页，每页10条）
func NewDefaultPage() *Page {
	return &Page{
		Page: 1,
		Size: 10,
	}
}
