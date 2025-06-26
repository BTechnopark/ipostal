package ipostal_api

type PageInfo struct {
	TotalItems  int  `json:"total_items,omitempty"`
	CurrentPage int  `json:"current_page"`
	TotalPages  int  `json:"total_pages,omitempty"`
	HasMore     bool `json:"has_more,omitempty"`
}

type ResponseData[T any] struct {
	Message  string    `json:"message"`
	Data     T         `json:"data,omitempty"`
	PageInfo *PageInfo `json:"page_info,omitempty"`
}

// SetError implements api_context.Response.
func (r *ResponseData[T]) SetError(err error) {
	r.Message = err.Error()
}
