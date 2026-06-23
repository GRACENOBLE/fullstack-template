package domain

// Page is a generic offset-based page of results.
type Page[T any] struct {
	Items    []T  `json:"items"`
	Total    int  `json:"total"`
	Page     int  `json:"page"`
	PageSize int  `json:"page_size"`
	HasMore  bool `json:"has_more"`
}

// CursorPage is a generic cursor-based page of results.
type CursorPage[T any] struct {
	Items      []T    `json:"items"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}

// PageRequest holds query-param pagination inputs.
type PageRequest struct {
	Page     int `form:"page"      binding:"min=0"`
	PageSize int `form:"page_size" binding:"min=0,max=100"`
}

// Defaults fills zero values with sensible defaults (page 1, size 20).
func (p *PageRequest) Defaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
}

// Offset returns the SQL OFFSET value for this page request.
// Normalises the page to at least 1 so the result is never negative.
func (p PageRequest) Offset() int {
	if p.Page <= 1 {
		return 0
	}
	return (p.Page - 1) * p.PageSize
}
