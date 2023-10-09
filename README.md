# Yzy Router

Creating easy routers for a Fiber framework, inspired by Rust attribute macros and Python decorators, based on AST parsing and reflection, just for fun :) Now with middlewares!

# [examples]

## Routers

### [[yzyrouter_test.go](yzyrouter_test.go)]

```go
// ..

app := fiber.New()
// ..

router := yzyrouter.New(app)
err := router.SetupControllers(
	// Basic controller
	controller.TaskController{},
	// Subdir controller
	index.IndexController{},
	// Diff package controller
	admin.AdminController{},
)
// ..
```

### [[tasks.go](example/tasks.go)]

```go
// ..
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
	return ctx.SendStatus(fiber.StatusOK)
}

// yzy:[Get:deleteTaskPath]
func (c TaskController) DeleteTask(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}

// yzy:[Get:deleteAllTasksPath]
func (c TaskController) DeleteAllTasks(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}
```

## Middlewares

### [[pong.go](/example/middleware/pong.go)]
```go
func New() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.WriteString("pong")
		return ctx.SendStatus(fiber.StatusOK)
	}
}
```

### [[yzyrouter_test.go](yzyrouter_test.go)]

```go
// ..

app := fiber.New()
// ..

router := yzyrouter.New(app)
router.RegisterMiddlewares(
	map[string]interface{}{
		"pong": pong.New(),
	},
)

err := router.SetupControllers(index.IndexController{}, admin.AdminController{})
// ..
```

### [[index.go](example/index/index.go)]

Here, pong is middleware on the `/index` route

```go
// ..
type IndexController struct{}

// yzy:[@pong]
// yzy:[Get:"/index"]
func (c IndexController) Index(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}

```

Here pong is middleware on all `AdminController` routes

```go
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
func (c AdminController) profilePath(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}
```
