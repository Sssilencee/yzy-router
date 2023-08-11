package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

const (
	apiRootPath = "/api/v1"

	tasksPath          = apiRootPath + "/tasks"
	addTaskPath        = tasksPath + "/add"
	deleteTaskPath     = tasksPath + "/delete"
	deleteAllTasksPath = deleteTaskPath + "All"
)

type TaskController struct{}

// yzyrouter[Get:tasksPath]
func (c TaskController) Tasks(ctx *fiber.Ctx) error {
	fmt.Println(ctx.OriginalURL())
	return nil
}

// yzyrouter[Get:addTaskPath]
func (c TaskController) AddTask(ctx *fiber.Ctx) error {
	fmt.Println(ctx.OriginalURL())
	return nil
}

// yzyrouter[Get:deleteTaskPath]
func (c TaskController) DeleteTask(ctx *fiber.Ctx) error {
	fmt.Println(ctx.OriginalURL())
	return nil
}

// yzyrouter[Get:deleteAllTasksPath]
func (c TaskController) DeleteAllTasks(ctx *fiber.Ctx) error {
	fmt.Println(ctx.OriginalURL())
	return nil
}
