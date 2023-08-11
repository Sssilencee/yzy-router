package admin

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type AdminController struct{}

// yzyrouter[Get:tasksPath]
func (c AdminController) Tasks(ctx *fiber.Ctx) error {
	fmt.Println(ctx.OriginalURL())
	return nil
}
