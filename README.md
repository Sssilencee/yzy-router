# Yzy Router

Creating easy routers for a Fiber framework, inspired by Rust attribute macros and Python decorators, based on AST parsing and reflection, just for fun :)

# [examples]

## [[yzyrouter_test.go](https://github.com/Sssilencee/yzyrouter/blob/main/yzyrouter.go)]

```go
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
	// ..

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

	// ..
}

```
## [[index.go](https://github.com/Sssilencee/yzyrouter/blob/main/example/index/index.go)]

```go
package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type IndexController struct{}

// yzy:[Get:"/index"]
func (c IndexController) Index(ctx *fiber.Ctx) error {
	fmt.Println(ctx.OriginalURL())
	return nil
}

```


## [[tasks.go](https://github.com/Sssilencee/yzyrouter/blob/main/example/tasks.go)]

```go
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

// yzy:[Get:addTaskPath]
func (c TaskController) AddTask(ctx *fiber.Ctx) error {
	fmt.Println(ctx.OriginalURL())
	return nil
}

// yzy:[Get:deleteTaskPath]
func (c TaskController) DeleteTask(ctx *fiber.Ctx) error {
	fmt.Println(ctx.OriginalURL())
	return nil
}

// yzy:[Get:deleteAllTasksPath]
func (c TaskController) DeleteAllTasks(ctx *fiber.Ctx) error {
	fmt.Println(ctx.OriginalURL())
	return nil
}
```
