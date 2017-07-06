package runtime

func NewDefaultContext() *Context {
	return &Context{
		Global: &JSObject{
			properties: map[string]Object{
				"console": console,
			},
		},
	}
}

type Context struct {
	Global *JSObject
}
