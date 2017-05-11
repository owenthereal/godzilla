package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/jingweno/godzilla/source"
)

type M map[string]interface{}

type MapUnmarshaler interface {
	UnmarshalMap(M)
}

type Compiler interface {
	Compile(*source.Code)
}

type Node interface {
	MapUnmarshaler
	Compiler
	fmt.Stringer
}

type Statement interface {
	Node
	statementNode()
}

type Declaration interface {
	Statement
	declarationNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Literal interface {
	Expression
	literalNode()
}

type Attr struct {
	Type  string
	Start int
	End   int
	Loc   *SourceLocation
}

func (a *Attr) UnmarshalMap(m M) {
	a.Type = convertString(m["type"])
	a.Start = convertInt(m["start"])
	a.End = convertInt(m["end"])
	a.Loc = unmarshalSourceLocation(convertMap(m["loc"]))
}

type SourceLocation struct {
	Start *Position
	End   *Position
}

func (s *SourceLocation) UnmarshalMap(m M) {
	s.Start = unmarshalPosition(convertMap(m["start"]))
	s.End = unmarshalPosition(convertMap(m["end"]))
}

type Position struct {
	Line   int `json:"line,omitempty"`
	Column int `json:"column,omitempty"`
}

func (p *Position) UnmarshalMap(m M) {
	p.Line = convertInt(m["line"])
	p.Column = convertInt(m["column"])
}

type File struct {
	*Attr
	Program *Program
}

func (f *File) UnmarshalMap(m M) {
	f.Attr = unmarshalAttr(m)
	f.Program = unmarshalProgram(convertMap(m["program"]))
}

func (f *File) Compile(code *source.Code) {
	f.Program.Compile(code)
}

func (f *File) String() string {
	return f.Program.String()
}

type Program struct {
	*Attr
	SourceType string
	Body       []Statement
}

func (p *Program) UnmarshalMap(m M) {
	p.Attr = unmarshalAttr(m)
	p.SourceType = convertString(m["sourceType"])
	p.Body = unmarshalStatements(convertSliceMap(m["body"]))
}

func (p *Program) Compile(code *source.Code) {
	for _, s := range p.Body {
		s.Compile(code)
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Body {
		out.WriteString(s.String())
	}

	return out.String()
}

type VariableDeclaration struct {
	*Attr
	Declarations []*VariableDeclarator
	Kind         string
}

func (v *VariableDeclaration) statementNode() {}

func (v *VariableDeclaration) declarationNode() {}

func (v *VariableDeclaration) UnmarshalMap(m M) {
	v.Attr = unmarshalAttr(m)
	v.Kind = convertString(m["kind"])
	v.Declarations = unmarshalVariableDeclarator(convertSliceMap(m["declarations"]))
}

func (v *VariableDeclaration) Compile(code *source.Code) {
	// TODO
}

func (v *VariableDeclaration) String() string {
	var out bytes.Buffer

	out.WriteString(v.Kind)
	for _, d := range v.Declarations {
		out.WriteString(d.String())
	}

	return out.String()
}

type VariableDeclarator struct {
	*Attr
	ID   *Identifier
	Init Expression
}

func (v *VariableDeclarator) UnmarshalMap(m M) {
	v.Attr = unmarshalAttr(m)
	v.Init = unmarshalExpression(convertMap(m["init"]))
	v.ID = unmarshalIdentifier(convertMap(m["id"]))
}

func (v *VariableDeclarator) Compile(code *source.Code) {
	// TODO
}

func (v *VariableDeclarator) String() string {
	var out bytes.Buffer

	out.WriteString(v.ID.String())
	out.WriteString(v.Init.String())

	return out.String()
}

type Identifier struct {
	*Attr
	Name string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) UnmarshalMap(m M) {
	i.Attr = unmarshalAttr(m)
	i.Name = convertString(m["name"])
}

func (i *Identifier) Compile(code *source.Code) {
	code.Write(strings.Title(i.Name))
}

func (i *Identifier) String() string {
	return i.Name
}

type StringLiteral struct {
	*Attr
	Extra *Extra
	Value string
}

func (s *StringLiteral) expressionNode() {}

func (s *StringLiteral) literalNode() {}

func (s *StringLiteral) UnmarshalMap(m M) {
	s.Attr = unmarshalAttr(m)
	s.Value = convertString(m["value"])
	s.Extra = unmarshalExtra(convertMap(m["extra"]))
}

func (s *StringLiteral) String() string {
	return fmt.Sprintf(`"%s"`, s.Value)
}

func (s *StringLiteral) Compile(code *source.Code) {
	code.Write(fmt.Sprintf(`"%s"`, s.Value))
}

type Extra struct {
	RawValue string
	Raw      string
}

func (e *Extra) UnmarshalMap(m M) {
	e.RawValue = convertString(m["rawValue"])
	e.Raw = convertString(m["raw"])
}

type CallExpression struct {
	*Attr
	Callee    Expression
	Arguments []Expression
}

func (c *CallExpression) expressionNode() {}

func (c *CallExpression) UnmarshalMap(m M) {
	c.Attr = unmarshalAttr(m)
	c.Callee = unmarshalExpression(convertMap(m["callee"]))
	c.Arguments = unmarshalExpressions(convertSliceMap(m["arguments"]))
}

func (c *CallExpression) Compile(code *source.Code) {
	c.Callee.Compile(code)
	code.Write("(")
	for i, arg := range c.Arguments {
		arg.Compile(code)
		if i != len(c.Arguments)-1 {
			code.Write(", ")
		}
	}
	code.Write(")\n")
}

func (c *CallExpression) String() string {
	var out bytes.Buffer

	out.WriteString(c.Callee.String())
	out.WriteString("(")

	var args []string
	for _, arg := range c.Arguments {
		args = append(args, arg.String())
	}
	out.WriteString(strings.Join(args, ", "))

	out.WriteString(")")

	return out.String()
}

type MemberExpression struct {
	*Attr
	Object   Expression
	Property Expression
	Computed bool
}

func (e *MemberExpression) expressionNode() {}

func (e *MemberExpression) UnmarshalMap(m M) {
	e.Attr = unmarshalAttr(m)
	e.Object = unmarshalExpression(convertMap(m["object"]))
	e.Property = unmarshalExpression(convertMap(m["property"]))
	e.Computed = convertBool(m["computed"])
}

func (e *MemberExpression) Compile(code *source.Code) {
	// TODO: ignoring computed value for now
	e.Object.Compile(code)
	code.Write(".")
	e.Property.Compile(code)
}

func (e *MemberExpression) String() string {
	return fmt.Sprintf("%s.%s", e.Object, e.Property)
}

type ExpressionStatement struct {
	*Attr
	Expression Expression
}

func (e *ExpressionStatement) statementNode() {}

func (e *ExpressionStatement) UnmarshalMap(m M) {
	e.Attr = unmarshalAttr(m)
	e.Expression = unmarshalExpression(convertMap(m["expression"]))
}

func (e *ExpressionStatement) Compile(code *source.Code) {
	e.Expression.Compile(code)
}

func (e *ExpressionStatement) String() string {
	return e.Expression.String()
}
