package yzyrouter

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type astPackage struct {
	*ast.Package
	set    *token.FileSet
	routes map[string]string
	mode   parserMode
}

type preamble struct {
	route          string
	method         string
	fn             string
	controllerName string
}

func pkgAst(path string, mode parserMode) (astPackage, error) {
	set := token.NewFileSet()
	pkg, err := parser.ParseDir(set, path, nil, parser.ParseComments)
	if err != nil {
		return astPackage{}, err
	}
	var p *ast.Package
	for k := range pkg {
		p = pkg[k]
	}
	return astPackage{
		Package: p,
		set:     set,
		mode:    mode,
	}, nil // only one pkg in this dir
}

func (ap astPackage) parsePreambles() []preamble {
	prs := make([]preamble, 0)

	for _, f := range ap.Files {
		for _, d := range f.Decls {
			if fn, isFn := d.(*ast.FuncDecl); isFn {

				preamble := fn.Doc.Text()
				if preamble == "" {
					continue
				}
				controllerName := fn.Recv.List[0].Type.(*ast.Ident).Name

				if p := strings.Split(preamble, "["); p[0] == "yzy:" {
					pr := ap.parsePreamble(p[1])
					pr.fn = fn.Name.String()
					pr.controllerName = controllerName

					prs = append(prs, pr)
				}
			}
		}
	}

	return prs
}

func (ap astPackage) parsePreamble(pr string) preamble {
	params := strings.Split(pr, ":")
	route := strings.TrimRight(params[1], "]\n") // ignore last `]`

	if string(route[0]) != "\"" || string(route[len(route)-1]) != "\"" {
		// Initialize map when meet first variable route
		if ap.routes == nil {
			routes := ap.parsePathVars()
			ap.routes = routes
		}
		route = ap.routes[route]
	} else {
		route = route[1 : len(route)-1]
	}
	return preamble{
		route:  route,
		method: strings.Title(strings.ToLower(params[0])), // method name before `:`
	}
}

func (ap astPackage) parsePathVars() map[string]string {
	routes := make(map[string]string)

	for _, f := range ap.Files {
		for _, d := range f.Decls {
			if gen, isGen := d.(*ast.GenDecl); isGen {
				if gen.Tok.String() == "const" || gen.Tok.String() == "var" {
					for _, s := range gen.Specs {
						if v, isV := s.(*ast.ValueSpec); isV {
							val, err := ap.parseAstExpr(v.Values[0])
							if err != nil && ap.mode == Debug {
								fmt.Println(err)
							}
							routes[v.Names[0].String()] = val
						}
					}
				}
			}
		}
	}

	return routes
}

func (ap astPackage) parseAstExpr(expr ast.Expr) (string, error) {
	switch e := expr.(type) {
	case *ast.Ident:
		return ap.parseAstExpr(e.Obj.Decl.(*ast.ValueSpec).Values[0])
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			return e.Value[1 : len(e.Value)-1], nil
		} else {
			return "", parserError{ap.set, e.Pos()}.exprTypeErr()
		}
	case *ast.BinaryExpr:
		op := e.Op.String()
		if op == "+" {
			xVal, xErr := ap.parseAstExpr(e.X)
			if xErr != nil {
				return "", xErr
			}
			yVal, yErr := ap.parseAstExpr(e.Y)
			if yErr != nil {
				return "", yErr
			}
			return xVal + yVal, nil
		} else {
			return "", parserError{ap.set, e.X.Pos()}.operatorErr(op)
		}
	default:
		return "", parserError{ap.set, e.Pos()}.undefinedNodeErr()
	}
}
