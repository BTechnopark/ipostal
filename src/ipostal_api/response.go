package ipostal_api

type ResponseData[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// SetError implements api_context.Response.
func (r *ResponseData[T]) SetError(err error) {
	r.Message = err.Error()
}
