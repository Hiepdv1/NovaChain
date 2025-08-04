package apperror

import (
	"ChainServer/internal/common/response"

	"github.com/gofiber/fiber/v2"
)

type AppError struct {
	Code    response.ErrorCode `json:"code"`
	Message string             `json:"message"`
	Err     error              `json:"-"`
	Status  uint               `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(message string, code response.ErrorCode, status uint, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
		Status:  status,
	}
}

func NotFound(msg string, err error) *AppError {
	return New(msg, response.ErrNotFound, fiber.StatusNotFound, err)
}

func Unauthorized(msg string, err error) *AppError {
	return New(msg, response.ErrUnauthorized, fiber.StatusUnauthorized, err)
}

func BadRequest(msg string, err error) *AppError {
	return New(msg, response.ErrBadRequest, fiber.StatusBadRequest, err)
}

func Internal(msg string, err error) *AppError {
	return New(msg, response.ErrInternal, fiber.StatusInternalServerError, err)
}
