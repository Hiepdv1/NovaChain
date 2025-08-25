package response

type ErrorType string

const (
	ErrBadRequest      ErrorType = "BAD_REQUEST"
	ErrUnauthorized    ErrorType = "UNAUTHORIZED"
	ErrNotFound        ErrorType = "NOT_FOUND"
	ErrInternal        ErrorType = "INTERNAL_SERVER_ERROR"
	ErrValidation      ErrorType = "VALIDATION_ERROR"
	ErrTooManyRequests ErrorType = "TOO_MANY_REQUESTS"
)
