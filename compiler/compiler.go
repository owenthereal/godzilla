package compiler

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jingweno/godzilla/ast"
	"github.com/jingweno/godzilla/source"
)

func Compile(f *ast.File) *source.Code {
	code := source.NewCode()

	c := &compiler{code}
	c.compile(f)

	return code
}

type compiler struct {
	code *source.Code
}

func (c *compiler) compile(f *ast.File) {
	c.compileProgram(f.Program)
}

func (c *compiler) compileProgram(p *ast.Program) {
	for _, s := range p.Body {
		c.compileStatement(s)
	}
}

// statements

func (c *compiler) compileStatement(s ast.Statement) {
	switch v := s.(type) {
	case *ast.ExpressionStatement:
		c.compileExpressionStatement(v)
	default:
		panic("unknown statement type " + getType(v))
	}
}

func (c *compiler) compileExpressionStatement(es *ast.ExpressionStatement) {
	c.compileExpression(es.Expression)
}

// expressions

func (c *compiler) compileExpression(e ast.Expression) {
	switch v := e.(type) {
	case *ast.CallExpression:
		c.compileCallExpression(v)
	case *ast.MemberExpression:
		c.compileMemberExpression(v)
	case *ast.Identifier:
		c.compileIdentifier(v)
	case *ast.StringLiteral:
		c.compileStringLiteral(v)
	default:
		panic("unknown expression type " + getType(v))
	}
}

func (c *compiler) compileCallExpression(ce *ast.CallExpression) {
	c.compileExpression(ce.Callee)
	c.code.Write("(")
	for i, arg := range ce.Arguments {
		c.compileExpression(arg)
		if i != len(ce.Arguments)-1 {
			c.code.Write(", ")
		}
	}
	c.code.Write(")\n")
}

// TODO: ignoring computed value for now
func (c *compiler) compileMemberExpression(me *ast.MemberExpression) {
	c.compileExpression(me.Object)
	c.code.Write(".")
	c.compileExpression(me.Property)
}

// TODO: look up from list of builtins
func (c *compiler) compileIdentifier(i *ast.Identifier) {
	c.code.Write(strings.Title(i.Name))
}

func (c *compiler) compileStringLiteral(s *ast.StringLiteral) {
	c.code.Write(fmt.Sprintf(`"%s"`, s.Value))
}

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
