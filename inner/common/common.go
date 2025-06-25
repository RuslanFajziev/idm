package common

import (
	"github.com/gofiber/fiber/v2"
)

type RequestValidationError struct {
	Message string
}

type AlreadyExistsError struct {
	Message string
}

type DbOperationError struct {
	Message string
}

type ResponseBody[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"error"`
	Data    T      `json:"data"`
}

func ErrResponse(
	c *fiber.Ctx,
	code int,
	message string,
) error {
	return c.Status(code).JSON(&ResponseBody[any]{
		Success: false,
		Message: message,
		Data:    nil,
	})
}

func OkResponse[T any](
	c *fiber.Ctx,
	data T,
) error {
	return c.JSON(&ResponseBody[T]{
		Success: true,
		Data:    data,
	})
}

func ResponseWithoutData(
	c *fiber.Ctx,
) error {
	return c.JSON(&ResponseBody[any]{
		Success: true,
	})
}

func (err RequestValidationError) Error() string {
	return err.Message
}

func (err AlreadyExistsError) Error() string {
	return err.Message
}

func (err DbOperationError) Error() string {
	return err.Message
}
