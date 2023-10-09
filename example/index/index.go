package controller

import (
	"github.com/gofiber/fiber/v2"
)

type IndexController struct{}

// yzy:[@pong]
// yzy:[Get:"/index"]
func (c IndexController) Index(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}
