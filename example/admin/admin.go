package admin

import (
	"github.com/gofiber/fiber/v2"
)

const (
	adminPath   = "/admin"
	profilePath = adminPath + "/profile"
)

// yzy:[@pong]
type AdminController struct{}

// yzy:[Get:adminPath]
func (c AdminController) Admin(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}

// yzy:[Get:profilePath]
func (c AdminController) Profile(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}
