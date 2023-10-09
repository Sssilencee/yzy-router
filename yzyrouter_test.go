package yzyrouter

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	controller "github.com/Sssilencee/yzyrouter/example"
	"github.com/Sssilencee/yzyrouter/example/admin"
	index "github.com/Sssilencee/yzyrouter/example/index"
	pong "github.com/Sssilencee/yzyrouter/example/middleware"
	"github.com/Sssilencee/yzyrouter/yzyrouter"
	"github.com/gofiber/fiber/v2"
)

const (
	serverUrl = "http://localhost:3000"

	indexPath   = "/index"
	adminPath   = "/admin"
	profilePath = adminPath + "/profile"

	apiRootPath = "/api/v1"

	tasksPath          = apiRootPath + "/tasks"
	addTaskPath        = tasksPath + "/add"
	deleteTaskPath     = tasksPath + "/delete"
	deleteAllTasksPath = deleteTaskPath + "All"
)

func TestSetupRouters(t *testing.T) {
	app := fiber.New()
	defer app.Shutdown()

	router := yzyrouter.New(app)
	err := router.SetupControllers(
		// Basic controller
		controller.TaskController{},
		// Subdir controller
		index.IndexController{},
		// Diff package controller
		admin.AdminController{},
	)
	if err != nil {
		t.Fatalf("setup controllers: %v", err)
	}

	{
		err := makeReq(app, serverUrl+indexPath)
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		err := makeReq(app, serverUrl+adminPath)
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		err := makeReq(app, serverUrl+addTaskPath)
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		err := makeReq(app, serverUrl+deleteTaskPath)
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		err := makeReq(app, serverUrl+deleteAllTasksPath)
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		err := makeReq(app, serverUrl+profilePath)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestSetupMiddlewares(t *testing.T) {
	app := fiber.New()
	defer app.Shutdown()

	router := yzyrouter.New(app)
	router.RegisterMiddlewares(
		map[string]interface{}{
			"pong": pong.New(),
		},
	)
	err := router.SetupControllers(index.IndexController{}, admin.AdminController{})
	if err != nil {
		t.Fatalf("setup controllers: %v", err)
	}

	{
		err := makeReqCheckPong(app, serverUrl+indexPath)
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		err := makeReqCheckPong(app, serverUrl+adminPath)
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		err := makeReqCheckPong(app, serverUrl+profilePath)
		if err != nil {
			t.Fatal(err)
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

func makeReqCheckPong(app *fiber.App, url string) error {
	const method = "GET"

	req := httptest.NewRequest(method, url, nil)
	res, err := app.Test(req)

	if err != nil {
		return fmt.Errorf("send req: %v", err)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read body: %v", err)
	}

	if string(b) != "pong" {
		return errors.New("invalid response")
	}

	return nil
}
