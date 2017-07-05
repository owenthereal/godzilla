package compiler

import (
	"fmt"
	"reflect"

	"github.com/jingweno/godzilla/ast"
	"github.com/jingweno/godzilla/runtime"
	"github.com/jingweno/godzilla/source"
)

func Compile(f *ast.File) *source.Code {
	code := source.NewCode()

	c := &compiler{
		code: code,
		ctx:  runtime.NewDefaultContext(),
	}
	c.compile(f)

	return code
}

type compiler struct {
	code *source.Code
	ctx  *runtime.Context
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
	case *ast.VariableDeclaration:
		c.compileVariableDeclaration(v)
	default:
		panic("unknown statement type " + getType(v))
	}
}

func (c *compiler) compileExpressionStatement(es *ast.ExpressionStatement) {
	c.compileExpression(es.Expression)
}

// TODO: ignore Kind for now
func (c *compiler) compileVariableDeclaration(vd *ast.VariableDeclaration) {
	for _, d := range vd.Declarations {
		c.compileVariableDeclarator(d)
	}
}

func (c *compiler) compileVariableDeclarator(vd *ast.VariableDeclarator) {
	c.writeLineNo(vd, vd.Attr)
	c.code.WriteLine(fmt.Sprintf("var %s Object", vd.ID))
	c.code.WriteLine(fmt.Sprintf("_ = %s", vd.ID))
	if vd.Init != nil {
		c.code.Write(fmt.Sprintf("%s = ", vd.ID))
		c.compileExpression(vd.Init)
		c.code.WriteLine("")
	}
	c.code.WriteLine(fmt.Sprintf(`global.DefineProperty("%s", %s)`, vd.ID, vd.ID))
}

func (c *compiler) writeLineNo(node ast.Node, attr *ast.Attr) {
	c.code.WriteLine(fmt.Sprintf(`// line %d: %s`, attr.Loc.Start.Line, node))
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
	c.writeLineNo(ce, ce.Attr)
	c.compileExpression(ce.Callee)
	c.code.Write("([]Object{")
	for i, arg := range ce.Arguments {
		c.compileExpression(arg)
		if i != len(ce.Arguments)-1 {
			c.code.Write(", ")
		}
	}
	c.code.Write("})\n")
}

// TODO: ignoring computed value for now
func (c *compiler) compileMemberExpression(me *ast.MemberExpression) {
	if me.Computed {
		panic("computed MemberExpression is not supported")
	}

	if builtInFunc := c.getBuiltinFunc(me.Object, me.Property); builtInFunc == "" {
		c.compileExpression(me.Object)
		c.code.Write(".")
		c.compileExpression(me.Property)
	} else {
		c.code.Write(builtInFunc)
	}
}

func (c *compiler) compileIdentifier(i *ast.Identifier) {
	c.code.Write(fmt.Sprintf(`global.GetProperty("%s")`, i.Name))
}

func (c *compiler) compileStringLiteral(s *ast.StringLiteral) {
	c.code.Write(fmt.Sprintf(`JSString("%s")`, s.Value))
}

func (c *compiler) getBuiltinFunc(objExp, propExp ast.Expression) string {
	oID, ok := objExp.(*ast.Identifier)
	if !ok {
		return ""
	}

	pID := propExp.(*ast.Identifier)
	if !ok {
		return ""
	}

	obj := c.ctx.Global.GetProperty(oID.Name)
	if obj == nil || obj.Type() != runtime.JS_OBJECT_TYPE_OBJECT {
		return ""
	}

	prop := (obj.(*runtime.JSObject)).GetProperty(pID.Name)
	if prop == nil || prop.Type() != runtime.JS_OBJECT_TYPE_FUNCTION {
		return ""
	}

	return (prop.(*runtime.JSFunction)).FuncName()
}

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
