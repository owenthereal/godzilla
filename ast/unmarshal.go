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
	e.RawValue = convertString(m["rawValue"])
	e.Raw = convertString(m["raw"])

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
	case "CallExpression":
		e = unmarshalCallExpression(m)
	case "MemberExpression":
		e = unmarshalMemberExpression(m)
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
