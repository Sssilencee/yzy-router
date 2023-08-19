package yzyrouter

import (
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type parserMode int

const (
	Silent parserMode = iota
	Debug
)

type yzyrouter struct {
	app  *fiber.App
	mode parserMode
}

func New(app *fiber.App, mode parserMode) yzyrouter {
	return yzyrouter{app, mode}
}

func (r yzyrouter) SetupControllers(controllers ...interface{}) {
	preambles := make(map[string][]preamble)
	for _, c := range controllers {
		path := reflect.ValueOf(c).Type().PkgPath()
		pkg, _ := pkgAst(strings.SplitN(path, "/", 4)[3], r.mode) // without mod dir

		// We don't want to parse one pkg for different controllers
		if _, exist := preambles[path]; !exist {
			preambles[path] = pkg.parsePreambles()
		}
		prs := preambles[path]

		controllerName := reflect.ValueOf(c).Type().Name()
		for _, p := range prs {
			if controllerName != p.controllerName {
				continue
			}
			f := reflect.ValueOf(c).MethodByName(p.fn)
			method := reflect.ValueOf(r.app).MethodByName(p.method)

			args := []reflect.Value{
				reflect.ValueOf(p.route),
				reflect.ValueOf(f.Interface().(func(*fiber.Ctx) error)),
			}
			method.Call(args)
		}
	}
}
