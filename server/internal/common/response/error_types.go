package response

type ErrorCode string

const (
	ErrBadRequest   ErrorCode = "BAD_REQUEST"
	ErrUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrNotFound     ErrorCode = "NOT_FOUND"
	ErrInternal     ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrValidation   ErrorCode = "VALIDATION_ERROR"
)
