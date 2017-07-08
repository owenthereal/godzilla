package ast

func unmarshalProgram(m m) *Program {
	p := &Program{}
	p.Attr = unmarshalAttr(m)
	p.SourceType = convertString(m["sourceType"])
	p.Body = unmarshalStatements(convertSliceMap(m["body"]))

	return p
}

func unmarshalAttr(m m) *Attr {
	a := &Attr{}
	a.Type = convertString(m["type"])
	a.Start = convertInt(m["start"])
	a.End = convertInt(m["end"])
	a.Loc = unmarshalSourceLocation(convertMap(m["loc"]))

	return a
}

func unmarshalSourceLocation(m m) *SourceLocation {
	sl := &SourceLocation{}
	sl.Start = unmarshalPosition(convertMap(m["start"]))
	sl.End = unmarshalPosition(convertMap(m["end"]))

	return sl
}

func unmarshalPosition(m m) *Position {
	p := &Position{}
	p.Line = convertInt(m["line"])
	p.Column = convertInt(m["column"])

	return p
}

func unmarshalExtra(m m) *Extra {
	e := &Extra{}
	e.RawValue = m["rawValue"]
	e.Raw = m["raw"]

	return e
}

// statements

func unmarshalStatements(m []m) []Statement {
	var s []Statement
	for _, mm := range m {
		s = append(s, unmarshalStatement(mm))
	}

	return s
}

func unmarshalStatement(m m) Statement {
	t := convertString(m["type"])
	var s Statement
	switch t {
	case "VariableDeclaration":
		s = unmarshalVariableDeclaration(m)
	case "ExpressionStatement":
		s = unmarshalExpressionStatement(m)
	case "BlockStatement":
		s = unmarshalBlockStatement(m)
	default:
		panic("unsupport statement type " + t)
	}

	return s
}

func unmarshalExpressionStatement(m m) *ExpressionStatement {
	e := &ExpressionStatement{}
	e.Attr = unmarshalAttr(m)
	e.Expression = unmarshalExpression(convertMap(m["expression"]))

	return e
}

func unmarshalVariableDeclaration(m m) *VariableDeclaration {
	v := &VariableDeclaration{}
	v.Attr = unmarshalAttr(m)
	v.Kind = convertString(m["kind"])
	v.Declarations = unmarshalVariableDeclarator(convertSliceMap(m["declarations"]))

	return v
}

// expressions

func unmarshalExpressions(m []m) []Expression {
	var e []Expression
	for _, mm := range m {
		e = append(e, unmarshalExpression(mm))
	}

	return e
}

func unmarshalExpression(m m) Expression {
	t := convertString(m["type"])
	var e Expression
	switch t {
	case "Identifier":
		e = unmarshalIdentifier(m)
	case "StringLiteral":
		e = unmarshalStringLiteral(m)
	case "NumericLiteral":
		e = unmarshalNumericLiteral(m)
	case "CallExpression":
		e = unmarshalCallExpression(m)
	case "MemberExpression":
		e = unmarshalMemberExpression(m)
	case "AssignmentExpression":
		e = unmarshalAssignmentExpression(m)
	case "BinaryExpression":
		e = unmarshalBinaryExpression(m)
	case "FunctionExpression":
		e = unmarshalFunctionExpression(m)
	default:
		panic("unsupport expression type " + t)
	}

	return e
}

func unmarshalIdentifier(m m) *Identifier {
	i := &Identifier{}
	i.Attr = unmarshalAttr(m)
	i.Name = convertString(m["name"])

	return i
}

func unmarshalCallExpression(m m) *CallExpression {
	c := &CallExpression{}
	c.Attr = unmarshalAttr(m)
	c.Callee = unmarshalExpression(convertMap(m["callee"]))
	c.Arguments = unmarshalExpressions(convertSliceMap(m["arguments"]))

	return c
}

func unmarshalMemberExpression(m m) *MemberExpression {
	e := &MemberExpression{}
	e.Attr = unmarshalAttr(m)
	e.Object = unmarshalExpression(convertMap(m["object"]))
	e.Property = unmarshalExpression(convertMap(m["property"]))
	e.Computed = convertBool(m["computed"])

	return e
}

func unmarshalAssignmentExpression(m m) *AssignmentExpression {
	a := &AssignmentExpression{}
	a.Attr = unmarshalAttr(m)
	a.Left = unmarshalExpression(convertMap(m["left"]))
	a.Right = unmarshalExpression(convertMap(m["right"]))
	a.Operator = AssignmentOperator(convertString(m["operator"]))

	return a
}

func unmarshalBinaryExpression(m m) *BinaryExpression {
	b := &BinaryExpression{}
	b.Attr = unmarshalAttr(m)
	b.Left = unmarshalExpression(convertMap(m["left"]))
	b.Right = unmarshalExpression(convertMap(m["right"]))
	b.Operator = BinaryOperator(convertString(m["operator"]))

	return b
}

func unmarshalFunctionExpression(m m) *FunctionExpression {
	f := &FunctionExpression{}
	f.Attr = unmarshalAttr(m)
	f.Function = unmarshalFunction(m)

	return f
}

func unmarshalFunction(m m) *Function {
	f := &Function{}
	if m["id"] != nil {
		f.ID = unmarshalIdentifier(convertMap(m["id"]))
	}
	f.Body = unmarshalBlockStatement(convertMap(m["body"]))
	f.Params = unmarshalFunctionParams(convertSliceMap(m["params"]))
	f.Generator = convertBool(m["generator"])
	f.Async = convertBool(m["async"])
	f.Expression = convertBool(m["expression"])

	return f
}

func unmarshalBlockStatement(m m) *BlockStatement {
	b := &BlockStatement{}
	b.Body = unmarshalStatements(convertSliceMap(m["body"]))
	// TODO: missing Directive

	return b
}

// TODO: only consider expression as params for now
func unmarshalFunctionParams(m []m) []Pattern {
	var p []Pattern
	for _, mm := range m {
		p = append(p, unmarshalExpression(mm))
	}

	return p
}

func unmarshalVariableDeclarator(m []m) []*VariableDeclarator {
	var d []*VariableDeclarator
	for _, mm := range m {
		dd := &VariableDeclarator{}
		dd.Attr = unmarshalAttr(mm)
		dd.ID = unmarshalIdentifier(convertMap(mm["id"]))
		if init := mm["init"]; init != nil {
			dd.Init = unmarshalExpression(convertMap(init))
		}

		d = append(d, dd)
	}

	return d
}

func unmarshalStringLiteral(m m) *StringLiteral {
	s := &StringLiteral{}
	s.Attr = unmarshalAttr(m)
	s.Value = convertString(m["value"])
	s.Extra = unmarshalExtra(convertMap(m["extra"]))

	return s
}

func unmarshalNumericLiteral(m m) *NumericLiteral {
	n := &NumericLiteral{}
	n.Attr = unmarshalAttr(m)
	n.Value = convertFloat(m["value"])
	n.Extra = unmarshalExtra(convertMap(m["extra"]))

	return n
}
