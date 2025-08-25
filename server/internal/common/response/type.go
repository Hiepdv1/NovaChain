package response

type ResponseBody struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
	Error      any    `json:"error,omitempty"`
	TraceID    string `json:"traceId,omitempty"`
}

type PaginationMeta struct {
	Limit       int  `json:"limit"`
	CurrentPage int  `json:"currentPage"`
	Total       int  `json:"total"`
	NextCursor  any  `json:"nextCursor"`
	HasMore     bool `json:"hasMore"`
}

type ListResponse struct {
	Success    bool           `json:"success"`
	StatusCode int            `json:"statusCode"`
	Message    string         `json:"message"`
	Data       any            `json:"data"`
	Meta       PaginationMeta `json:"meta"`
	TraceID    string         `json:"traceId"`
}
