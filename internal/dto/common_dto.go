package dto

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ErrorResponse `json:"errors"`
}

// SetDefaultPagination sets default values for pagination
func (q *NoteQueryRequest) SetDefaults() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}

	if q.SortBy == "" {
		q.SortBy = "created_at"
	}
	if q.SortDir == "" {
		q.SortDir = "desc"
	}
}

// CalculateTotalPages calculates total pages from total items
func CalculateTotalPages(totalItems int64, pageSize int) int {
	if pageSize == 0 {
		return 0
	}
	pages := int(totalItems) / pageSize
	if int(totalItems)%pageSize > 0 {
		pages++
	}
	return pages

}
