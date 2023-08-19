package controller

import (
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

// yzy:[Get:addTaskPath]
func (c TaskController) AddTask(ctx *fiber.Ctx) error {
	ctx.SendStatus(fiber.StatusOK)
	return nil
}

// yzy:[Get:deleteTaskPath]
func (c TaskController) DeleteTask(ctx *fiber.Ctx) error {
	ctx.SendStatus(fiber.StatusOK)
	return nil
}

// yzy:[Get:deleteAllTasksPath]
func (c TaskController) DeleteAllTasks(ctx *fiber.Ctx) error {
	ctx.SendStatus(fiber.StatusOK)
	return nil
}
