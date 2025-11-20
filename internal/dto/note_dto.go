package dto

// CreateNoteRequest represents note creation request
type CreateNoteRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=255"`
	Content string `json:"content" validate:"max=10000"`
}

// UpdateNoteRequest represents note update request
type UpdateNoteRequest struct {
	Title   string `json:"title" validate:"omitempty,min=1,max=255"`
	Content string `json:"content" validate:"omitempty,max=10000"`
}

// NoteQueryRequest represents note list query parameters
type NoteQueryRequest struct {
	Page     int    `query:"page" validate:"omitempty,min=1"`
	PageSize int    `query:"page_size" validate:"omitempty,min=1,max=100"`
	Search   string `query:"search" validate:"omitempty,max=100"`
	SortBy   string `query:"sort_by" validate:"omitempty,oneof=created_at updated_at title"`
	SortDir  string `query:"sort_dir" validate:"omitempty,oneof=asc desc"`
}

// NoteResponse represents note data in response
type NoteResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// NoteListResponse represents paginated note list response
type NoteListResponse struct {
	Notes      []NoteResponse `json:"notes"`
	Pagination PaginationMeta `json:"Pagination"`
}
