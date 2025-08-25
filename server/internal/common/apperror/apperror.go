package apperror

import (
	"ChainServer/internal/common/response"

	"github.com/gofiber/fiber/v2"
)

type AppError struct {
	ErrType response.ErrorType `json:"errorType"`
	Message string             `json:"message"`
	Err     error              `json:"-"`
	Status  uint               `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(message string, errType response.ErrorType, status uint, err error) *AppError {
	return &AppError{
		ErrType: errType,
		Message: message,
		Err:     err,
		Status:  status,
	}
}

func (apperr *AppError) Reponse(c *fiber.Ctx) error {
	traceID := response.GetTraceID(c)

	res := response.ResponseBody{
		Success:    apperr.Status >= 200 && apperr.Status < 300,
		StatusCode: int(apperr.Status),
		Message:    apperr.Message,
		TraceID:    traceID,
		Error:      apperr.Err,
	}

	return c.Status(int(apperr.Status)).JSON(res)

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

func TooManyRequests(msg string, err error) *AppError {
	return New(msg, response.ErrTooManyRequests, fiber.StatusTooManyRequests, err)
}
