package runtime

import "fmt"

var (
	console = &JSObject{
		properties: map[string]Object{
			"log": &JSFunction{
				fn: Console_Log,
			},
		},
	}
)

func Console_Log(data []Object) {
	var i []interface{}
	for _, d := range data {
		i = append(i, d)
	}

	fmt.Println(i...)
}
