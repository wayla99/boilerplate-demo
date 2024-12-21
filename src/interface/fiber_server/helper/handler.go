package helper

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

const OK string = "OK"

type ErrorResponse struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}

type ErrorMapInfo struct {
	StatusCode int
	ErrorCode  string
}

func SendError(c *fiber.Ctx, status int, err error, errCode string) error {
	c.Locals("error", err)

	return c.Status(status).JSON(ErrorResponse{
		Error:     err.Error(),
		ErrorCode: errCode,
	})
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	unwrapErr := errors.Unwrap(err)
	if unwrapErr == nil {
		unwrapErr = err
	}

	return SendError(c, 500, err, "unexpected_error")
}
