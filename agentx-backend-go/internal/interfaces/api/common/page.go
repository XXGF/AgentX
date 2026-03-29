package common

// PageRequest 分页请求参数（对应 Java 的 Page 基类）
type PageRequest struct {
	Page     int `form:"page" json:"page"`         // 页码，从1开始
	PageSize int `form:"pageSize" json:"pageSize"`  // 每页大小
}

// GetPage 获取页码（默认1）
func (p *PageRequest) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetPageSize 获取每页大小（默认15，最大100）
func (p *PageRequest) GetPageSize() int {
	if p.PageSize <= 0 {
		return 15
	}
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}

// GetOffset 获取数据库查询偏移量
func (p *PageRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}

// PageResponse 分页响应（对应 Java 的 MyBatis-Plus Page 返回结构）
type PageResponse[T any] struct {
	Records []T   `json:"records"` // 数据记录
	Total   int64 `json:"total"`   // 总记录数
	Size    int   `json:"size"`    // 每页大小
	Current int   `json:"current"` // 当前页码
	Pages   int64 `json:"pages"`   // 总页数
}

// NewPageResponse 创建分页响应
func NewPageResponse[T any](records []T, total int64, page, pageSize int) *PageResponse[T] {
	pages := total / int64(pageSize)
	if total%int64(pageSize) > 0 {
		pages++
	}
	return &PageResponse[T]{
		Records: records,
		Total:   total,
		Size:    pageSize,
		Current: page,
		Pages:   pages,
	}
}
