package handlers

import (
	"github.com/gofiber/fiber/v2"
	"news/shared/failure"
)

var SuccessResponse = func(ctx *fiber.Ctx, status int, data interface{}) error {
	return ctx.Status(status).JSON(&fiber.Map{
		"status": status,
		"data":   data,
	})
}

var ErrorResponse = func(ctx *fiber.Ctx, err error) error {
	code := failure.GetCode(err)
	return ctx.Status(code).JSON(fiber.Map{
		"code":  code,
		"error": err.Error(),
	})
}
