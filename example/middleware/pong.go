package pong

import (
	"github.com/gofiber/fiber/v2"
)

func New() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.WriteString("pong")
		return ctx.SendStatus(fiber.StatusOK)
	}
}
