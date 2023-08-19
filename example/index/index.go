package controller

import (
	"github.com/gofiber/fiber/v2"
)

type IndexController struct{}

// yzy:[Get:"/index"]
func (c IndexController) Index(ctx *fiber.Ctx) error {
	ctx.SendStatus(fiber.StatusOK)
	return nil
}
