package yzyrouter

import (
	"fmt"
	"go/token"
)

const (
	opErrMsg            = "only concatianation operator is allowed not \"%s\": %s"
	undefinedNodeErrMsg = "find undefined ast node: %s"
	exprTypeErrMsg      = "expr isn't a type of string: %s"
)

type parserError struct {
	set *token.FileSet
	pos token.Pos
}

func (e parserError) operatorErr(op string) error {
	return fmt.Errorf(opErrMsg, op, e.set.Position(e.pos).String())
}

func (e parserError) undefinedNodeErr() error {
	return fmt.Errorf(undefinedNodeErrMsg, e.set.Position(e.pos).String())
}

func (e parserError) exprTypeErr() error {
	return fmt.Errorf(exprTypeErrMsg, e.set.Position(e.pos).String())
}
