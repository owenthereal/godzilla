package ast

func unmarshalProgram(m M) *Program {
	p := &Program{}
	p.UnmarshalMap(m)

	return p
}

func unmarshalAttr(m M) *Attr {
	attr := &Attr{}
	attr.UnmarshalMap(m)

	return attr
}

func unmarshalSourceLocation(m M) *SourceLocation {
	sl := &SourceLocation{}
	sl.UnmarshalMap(m)

	return sl
}

func unmarshalPosition(m M) *Position {
	p := &Position{}
	p.UnmarshalMap(m)

	return p
}

func unmarshalIdentifier(m M) *Identifier {
	i := &Identifier{}
	i.UnmarshalMap(convertMap(m["id"]))

	return i
}

func unmarshalExtra(m M) *Extra {
	e := &Extra{}
	e.UnmarshalMap(m)

	return e
}

func unmarshalDeclaration(m M) Declaration {
	t := convertString(m["type"])
	var d Declaration
	switch t {
	case "VariableDeclaration":
		d = &VariableDeclaration{}
	default:
		panic("unsupport declaration type " + t)
	}

	d.UnmarshalMap(m)

	return d
}

func unmarshalExpressions(m []M) []Expression {
	var e []Expression
	for _, mm := range m {
		e = append(e, unmarshalExpression(mm))
	}

	return e
}

func unmarshalExpression(m M) Expression {
	t := convertString(m["type"])
	var e Expression
	switch t {
	case "Identifier":
		e = &Identifier{}
	case "StringLiteral":
		e = &StringLiteral{}
	case "CallExpression":
		e = &CallExpression{}
	case "MemberExpression":
		e = &MemberExpression{}
	default:
		panic("unsupport expression type " + t)
	}

	e.UnmarshalMap(m)

	return e
}

func unmarshalVariableDeclarator(m []M) []*VariableDeclarator {
	var d []*VariableDeclarator
	for _, mm := range m {
		dd := &VariableDeclarator{}
		dd.UnmarshalMap(mm)
		d = append(d, dd)
	}

	return d
}

func unmarshalStatements(m []M) []Statement {
	var s []Statement
	for _, mm := range m {
		s = append(s, unmarshalStatement(mm))
	}

	return s
}

func unmarshalStatement(m M) Statement {
	t := convertString(m["type"])
	var s Statement
	switch t {
	case "VariableDeclaration":
		s = &VariableDeclaration{}
	case "ExpressionStatement":
		s = &ExpressionStatement{}
	default:
		panic("unsupport statement type " + t)
	}

	s.UnmarshalMap(m)

	return s
}

func convertSliceMap(i interface{}) []M {
	var m []M
	for _, mm := range i.([]interface{}) {
		m = append(m, convertMap(mm))
	}

	return m
}

func convertMap(i interface{}) M {
	return M(i.(map[string]interface{}))
}

func convertString(i interface{}) string {
	return i.(string)
}

func convertInt(i interface{}) int {
	return int(i.(float64))
}

func convertBool(i interface{}) bool {
	return i.(bool)
}
