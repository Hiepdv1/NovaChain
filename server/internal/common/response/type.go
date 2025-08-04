package response

type ResponseBody struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
	Error      any    `json:"error,omitempty"`
	TraceID    string `json:"traceId,omitempty"`
	Stack      any    `json:"stack,omitempty"`
}

type PaginationMeta struct {
	Limit       int  `json:"limit"`
	CurrentPage int  `json:"currentPage"`
	Total       int  `json:"total"`
	NextCursor  any  `json:"nextCursor"`
	HasMore     bool `json:"hasMore"`
}

type ListResponse struct {
	StatusCode int            `json:"statusCode"`
	Message    string         `json:"message"`
	Data       any            `json:"data"`
	Meta       PaginationMeta `json:"meta"`
	TraceID    string         `json:"traceId"`
}
