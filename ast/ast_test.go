package ast

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestUnmarshalJSON(t *testing.T) {
	s := `{"type":"File","start":0,"end":31,"loc":{"start":{"line":1,"column":0},"end":{"line":2,"column":0}},"program":{"type":"Program","start":0,"end":31,"loc":{"start":{"line":1,"column":0},"end":{"line":2,"column":0}},"sourceType":"script","body":[{"type":"ExpressionStatement","start":0,"end":30,"loc":{"start":{"line":1,"column":0},"end":{"line":1,"column":30}},"expression":{"type":"CallExpression","start":0,"end":30,"loc":{"start":{"line":1,"column":0},"end":{"line":1,"column":30}},"callee":{"type":"MemberExpression","start":0,"end":11,"loc":{"start":{"line":1,"column":0},"end":{"line":1,"column":11}},"object":{"type":"Identifier","start":0,"end":7,"loc":{"start":{"line":1,"column":0},"end":{"line":1,"column":7},"identifierName":"console"},"name":"console"},"property":{"type":"Identifier","start":8,"end":11,"loc":{"start":{"line":1,"column":8},"end":{"line":1,"column":11},"identifierName":"log"},"name":"log"},"computed":false},"arguments":[{"type":"StringLiteral","start":12,"end":29,"loc":{"start":{"line":1,"column":12},"end":{"line":1,"column":29}},"extra":{"rawValue":"Hello, Godzilla","raw":"'Hello, Godzilla'"},"value":"Hello, Godzilla"}]}}],"directives":[]},"comments":[],"tokens":[{"type":{"label":"name","beforeExpr":false,"startsExpr":true,"rightAssociative":false,"isLoop":false,"isAssign":false,"prefix":false,"postfix":false,"binop":null},"value":"console","start":0,"end":7,"loc":{"start":{"line":1,"column":0},"end":{"line":1,"column":7}}},{"type":{"label":".","beforeExpr":false,"startsExpr":false,"rightAssociative":false,"isLoop":false,"isAssign":false,"prefix":false,"postfix":false,"binop":null,"updateContext":null},"start":7,"end":8,"loc":{"start":{"line":1,"column":7},"end":{"line":1,"column":8}}},{"type":{"label":"name","beforeExpr":false,"startsExpr":true,"rightAssociative":false,"isLoop":false,"isAssign":false,"prefix":false,"postfix":false,"binop":null},"value":"log","start":8,"end":11,"loc":{"start":{"line":1,"column":8},"end":{"line":1,"column":11}}},{"type":{"label":"(","beforeExpr":true,"startsExpr":true,"rightAssociative":false,"isLoop":false,"isAssign":false,"prefix":false,"postfix":false,"binop":null},"start":11,"end":12,"loc":{"start":{"line":1,"column":11},"end":{"line":1,"column":12}}},{"type":{"label":"string","beforeExpr":false,"startsExpr":true,"rightAssociative":false,"isLoop":false,"isAssign":false,"prefix":false,"postfix":false,"binop":null,"updateContext":null},"value":"Hello, Godzilla","start":12,"end":29,"loc":{"start":{"line":1,"column":12},"end":{"line":1,"column":29}}},{"type":{"label":")","beforeExpr":false,"startsExpr":false,"rightAssociative":false,"isLoop":false,"isAssign":false,"prefix":false,"postfix":false,"binop":null},"start":29,"end":30,"loc":{"start":{"line":1,"column":29},"end":{"line":1,"column":30}}},{"type":{"label":"eof","beforeExpr":false,"startsExpr":false,"rightAssociative":false,"isLoop":false,"isAssign":false,"prefix":false,"postfix":false,"binop":null,"updateContext":null},"start":31,"end":31,"loc":{"start":{"line":2,"column":0},"end":{"line":2,"column":0}}}]}`
	got := &File{}
	if err := json.Unmarshal([]byte(s), got); err != nil {
		t.Fatalf("json unmarshal has error: %s", err)
	}

	want := &File{
		Attr: &Attr{
			Type:  "File",
			Start: 0,
			End:   31,
			Loc: &SourceLocation{
				Start: &Position{
					Line:   1,
					Column: 0,
				},
				End: &Position{
					Line:   2,
					Column: 0,
				},
			},
		},
	}

	if !reflect.DeepEqual(want.Attr, got.Attr) {
		t.Fatalf("file not equal: want=%s got=%s", want, got)
	}
}
