package main

import (
	"testing"
	controller "yzyrouter/example"
	"yzyrouter/example/admin"
	index "yzyrouter/example/index"
	"yzyrouter/yzyrouter"

	"github.com/gofiber/fiber/v2"
)

func TestSetup(t *testing.T) {
	app := fiber.New()
	router := yzyrouter.New(app, yzyrouter.Debug)
	router.SetupControllers(
		// Basic controller
		controller.TaskController{},
		// Subdir controller
		index.IndexController{},
		// Diff package controller
		admin.AdminController{},
	)
	app.Listen(":3000")
}
