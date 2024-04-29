package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func NewTestApp() *fiber.App {
	return fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if apiError, ok := err.(Error); ok {
				return c.Status(apiError.Code).JSON(apiError)
			}
			internalError := NewError(http.StatusInternalServerError, "internal server error")
			return c.Status(internalError.Code).JSON(internalError)
		}})
}
