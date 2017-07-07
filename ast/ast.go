package ast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type m map[string]interface{}

type Node interface {
	fmt.Stringer
	GetAttr() *Attr
}

type File struct {
	*Attr
	Program *Program
}

func (f *File) UnmarshalJSON(data []byte) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	f.Attr = unmarshalAttr(m)
	f.Program = unmarshalProgram(convertMap(m["program"]))

	return nil
}

func (f *File) GetAttr() *Attr {
	return f.Attr
}

func (f *File) String() string {
	if f.Program == nil {
		return ""
	}

	return f.Program.String()
}

type Program struct {
	*Attr
	SourceType string
	Body       []Statement
}

func (p *Program) GetAttr() *Attr {
	return p.Attr
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Body {
		out.WriteString(s.String())
	}

	return out.String()
}

type Attr struct {
	Type  string
	Start int
	End   int
	Loc   *SourceLocation
}

type SourceLocation struct {
	Start *Position
	End   *Position
}

type Position struct {
	Line   int
	Column int
}

type Extra struct {
	RawValue interface{}
	Raw      interface{}
}

// statements

type Statement interface {
	Node
	statementNode()
}

type ExpressionStatement struct {
	*Attr
	Expression Expression
}

func (e *ExpressionStatement) statementNode() {}

func (e *ExpressionStatement) GetAttr() *Attr {
	return e.Attr
}

func (e *ExpressionStatement) String() string {
	return e.Expression.String()
}

// declarations

type Declaration interface {
	Statement
	declarationNode()
}

type VariableDeclaration struct {
	*Attr
	Declarations []*VariableDeclarator
	Kind         string
}

func (v *VariableDeclaration) statementNode() {}

func (v *VariableDeclaration) declarationNode() {}

func (v *VariableDeclaration) GetAttr() *Attr {
	return v.Attr
}

func (v *VariableDeclaration) String() string {
	var out bytes.Buffer

	out.WriteString(v.Kind)
	out.WriteString(" ")
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

func (v *VariableDeclarator) GetAttr() *Attr {
	return v.Attr
}

func (v *VariableDeclarator) String() string {
	var out bytes.Buffer

	out.WriteString(v.ID.String())
	if v.Init != nil {
		out.WriteString(" = ")
		out.WriteString(v.Init.String())
	}

	return out.String()
}

// expressions

type Expression interface {
	Node
	expressionNode()
}

type Identifier struct {
	*Attr
	Name string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) GetAttr() *Attr {
	return i.Attr
}

func (i *Identifier) String() string {
	return i.Name
}

type CallExpression struct {
	*Attr
	Callee    Expression
	Arguments []Expression
}

func (c *CallExpression) expressionNode() {}

func (c *CallExpression) GetAttr() *Attr {
	return c.Attr
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

func (m *MemberExpression) GetAttr() *Attr {
	return m.Attr
}

func (e *MemberExpression) String() string {
	return fmt.Sprintf("%s.%s", e.Object, e.Property)
}

type AssignmentExpression struct {
	*Attr
	Operator AssignmentOperator
	Left     Expression
	Right    Expression
}

func (a *AssignmentExpression) expressionNode() {}

func (a *AssignmentExpression) GetAttr() *Attr {
	return a.Attr
}

func (a *AssignmentExpression) String() string {
	return fmt.Sprintf("%s %s %s", a.Left, a.Operator, a.Right)
}

type AssignmentOperator string

type BinaryExpression struct {
	*Attr
	Operator BinaryOperator
	Left     Expression
	Right    Expression
}

func (b *BinaryExpression) expressionNode() {}

func (b *BinaryExpression) GetAttr() *Attr {
	return b.Attr
}

func (a *BinaryExpression) String() string {
	return fmt.Sprintf("%s %s %s", a.Left, a.Operator, a.Right)
}

type BinaryOperator string

// literals

type Literal interface {
	Expression
	literalNode()
}

type StringLiteral struct {
	*Attr
	Extra *Extra
	Value string
}

func (s *StringLiteral) expressionNode() {}

func (s *StringLiteral) literalNode() {}

func (s *StringLiteral) GetAttr() *Attr {
	return s.Attr
}

func (s *StringLiteral) String() string {
	return fmt.Sprintf(`"%s"`, s.Value)
}

// TODO: Value is always float64
// Can delay conversion and adapt to int vs. float
type NumericLiteral struct {
	*Attr
	Extra *Extra
	Value float64
}

func (n *NumericLiteral) expressionNode() {}

func (n *NumericLiteral) literalNode() {}

func (n *NumericLiteral) GetAttr() *Attr {
	return n.Attr
}

func (n *NumericLiteral) String() string {
	return fmt.Sprintf("%.16f", n.Value)
}
