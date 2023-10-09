package yzyrouter

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Yzyrouter struct {
	app         *fiber.App
	middlewares map[string]interface{}
}

func New(app *fiber.App) Yzyrouter {
	return Yzyrouter{app: app}
}

func NewWithMiddleware(app *fiber.App, middleware map[string]interface{}) Yzyrouter {
	return Yzyrouter{app, middleware}
}

func (r Yzyrouter) SetupControllers(controllers ...interface{}) error {
	preambles := make(map[string][]preamble)
	for _, c := range controllers {
		path := reflect.ValueOf(c).Type().PkgPath()
		n, i := 2, 1
		// Only works with github.com/<USERNAME>/<MODULE_NAME>
		if strings.HasPrefix(path, "github.com") {
			n, i = 4, 3
		}
		pkg, _ := pkgAst(strings.SplitN(path, "/", n)[i]) // without mod dir

		// We don't want to parse one pkg for different controllers
		if _, exist := preambles[path]; !exist {
			prs, err := pkg.parsePreambles()
			if err != nil {
				return err
			}
			preambles[path] = prs
		}
		prs := preambles[path]

		controllerName := reflect.ValueOf(c).Type().Name()
		for _, p := range prs {
			if controllerName != p.controllerName {
				continue
			}

			if len(p.middlewares) != 0 {
				args := []interface{}{p.route}
				for _, name := range p.middlewares {
					f, exist := r.middlewares[name]
					if !exist {
						return fmt.Errorf("invalid middleware name: %s", name)
					}
					args = append(args, f)
				}

				r.app.Use(args...)
			}

			f := reflect.ValueOf(c).MethodByName(p.fn)
			if f.IsValid() == false {
				return fmt.Errorf("\"%s\" method \"%s\" probably unexported", controllerName, p.fn)
			}

			method := reflect.ValueOf(r.app).MethodByName(p.method)

			args := []reflect.Value{
				reflect.ValueOf(p.route),
				reflect.ValueOf(f.Interface().(func(*fiber.Ctx) error)),
			}
			method.Call(args)
		}
	}

	return nil
}

func (r *Yzyrouter) RegisterMiddlewares(middlewares map[string]interface{}) {
	r.middlewares = middlewares
}
