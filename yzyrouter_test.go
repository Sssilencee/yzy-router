package yzyrouter

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	controller "github.com/Sssilencee/yzyrouter/example"
	"github.com/Sssilencee/yzyrouter/example/admin"
	index "github.com/Sssilencee/yzyrouter/example/index"
	"github.com/Sssilencee/yzyrouter/yzyrouter"
	"github.com/gofiber/fiber/v2"
)

const (
	serverUrl = "http://localhost:3000"

	indexPath = "/index"
	adminPath = "/admin"

	apiRootPath = "/api/v1"

	tasksPath          = apiRootPath + "/tasks"
	addTaskPath        = tasksPath + "/add"
	deleteTaskPath     = tasksPath + "/delete"
	deleteAllTasksPath = deleteTaskPath + "All"
)

func TestSetup(t *testing.T) {
	app := fiber.New()
	defer app.Shutdown()

	router := yzyrouter.New(app, yzyrouter.Debug)
	router.SetupControllers(
		// Basic controller
		controller.TaskController{},
		// Subdir controller
		index.IndexController{},
		// Diff package controller
		admin.AdminController{},
	)

	{
		err := makeReq(app, serverUrl+indexPath)
		if err != nil {
			log.Fatal(err)
		}
	}
	{
		err := makeReq(app, serverUrl+adminPath)
		if err != nil {
			log.Fatal(err)
		}
	}
	{
		err := makeReq(app, serverUrl+addTaskPath)
		if err != nil {
			log.Fatal(err)
		}
	}
	{
		err := makeReq(app, serverUrl+deleteTaskPath)
		if err != nil {
			log.Fatal(err)
		}
	}
	{
		err := makeReq(app, serverUrl+deleteAllTasksPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func makeReq(app *fiber.App, url string) error {
	const method = "GET"

	req := httptest.NewRequest(method, url, nil)
	res, err := app.Test(req)

	if err != nil {
		return fmt.Errorf("index controller error: %s", err.Error())
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("index controller status: %s", res.Status)
	}

	return nil
}
