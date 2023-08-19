package admin

import (
	"github.com/gofiber/fiber/v2"
)

const adminPath = "/admin"

type AdminController struct{}

// yzy:[Get:adminPath]
func (c AdminController) Tasks(ctx *fiber.Ctx) error {
	ctx.SendStatus(fiber.StatusOK)
	return nil
}
