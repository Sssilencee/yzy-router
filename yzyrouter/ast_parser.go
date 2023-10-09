package yzyrouter

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"strings"
)

type astPackage struct {
	*ast.Package
	set    *token.FileSet
	routes map[string]string
}

type preamble struct {
	route          string
	method         string
	fn             string
	controllerName string
	middlewares    []string
}

func pkgAst(path string) (astPackage, error) {
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
	}, nil // only one pkg in this dir
}

func (ap astPackage) parsePreambles() ([]preamble, error) {
	// `if len(preamble) < 8` because: yzy:[@_] - 8 symbols

	prs := make([]preamble, 0)
	for _, f := range ap.Files {
		for _, d := range f.Decls {
			if fn, isFn := d.(*ast.FuncDecl); isFn {

				preamble := fn.Doc.Text()
				if len(preamble) < 8 {
					continue
				}

				middlewares := make([]string, 0)
				controllerName := fn.Recv.List[0].Type.(*ast.Ident).Name
				preambles := strings.Split(preamble, "\n")
				for _, p := range preambles {
					if len(p) < 8 || p[:4] != "yzy:" {
						break
					}

					if p[5] == '@' {
						middlewares = append(middlewares, p[6:len(p)-1])
					} else {
						pr, err := ap.parsePreamble(p[5:])
						if err != nil {
							return nil, err
						}
						pr.fn = fn.Name.String()
						pr.controllerName = controllerName
						pr.middlewares = middlewares

						prs = append(prs, pr)
					}
				}
			}
		}
	}

	mws := make(map[string][]string, 0)
	pkg := doc.New(ap.Package, "./", doc.AllMethods)
	for _, t := range pkg.Types {
		if len(t.Doc) < 8 {
			continue
		}

		preambles := strings.Split(t.Doc, "\n")
		for _, p := range preambles {
			if len(p) < 8 || p[:4] != "yzy:" || p[5] != '@' {
				continue
			}

			_, exist := mws[t.Name]
			if !exist {
				mws[t.Name] = make([]string, 0, 1)
			}
			mws[t.Name] = append(mws[t.Name], p[6:len(p)-1])
		}
	}

	for i, p := range prs {
		m, exist := mws[p.controllerName]
		if exist {
			prs[i].middlewares = append(p.middlewares, m...)
		}
	}

	return prs, nil
}

func (ap astPackage) parsePreamble(pr string) (preamble, error) {
	params := strings.Split(pr, ":")
	route := strings.TrimRight(params[1], "]") // ignore last `]`

	if string(route[0]) != "\"" || string(route[len(route)-1]) != "\"" {
		// Initialize map when meet first variable route
		if ap.routes == nil {
			routes, err := ap.parsePathVars()
			if err != nil {
				return preamble{}, err
			}
			ap.routes = routes
		}
		route = ap.routes[route]
	} else {
		route = route[1 : len(route)-1]
	}
	return preamble{
		route:  route,
		method: strings.Title(strings.ToLower(params[0])), // method name before `:`
	}, nil
}

func (ap astPackage) parsePathVars() (map[string]string, error) {
	routes := make(map[string]string)

	for _, f := range ap.Files {
		for _, d := range f.Decls {
			if gen, isGen := d.(*ast.GenDecl); isGen {
				if gen.Tok.String() == "const" || gen.Tok.String() == "var" {
					for _, s := range gen.Specs {
						if v, isV := s.(*ast.ValueSpec); isV {
							val, err := ap.parseAstExpr(v.Values[0])
							if err != nil {
								return nil, fmt.Errorf("ast parsing error: %v", err)
							}
							routes[v.Names[0].String()] = val
						}
					}
				}
			}
		}
	}

	return routes, nil
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
