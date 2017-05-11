package source

import (
	"bytes"
	"io"
	"strings"
	"text/template"
)

const tmpl = `package main

import (
	. "github.com/jingweno/godzilla/runtime"
)

func main() {
	{{.}}
}`

func NewCode() *Code {
	return &Code{bytes.NewBuffer(nil)}
}

type Code struct {
	buf *bytes.Buffer
}

func (c *Code) WriteTo(w io.Writer) error {
	t, err := template.New("main").Parse(tmpl)
	if err != nil {
		return err
	}

	return t.Execute(w, strings.TrimSpace(c.buf.String()))
}

func (c *Code) String() string {
	result := bytes.NewBuffer(nil)
	err := c.WriteTo(result)
	if err != nil {
		panic(err)
	}

	return result.String()
}

func (c *Code) Write(s string) {
	c.buf.WriteString(s)
}

func (c *Code) WriteLine(s string) {
	c.Write(s)
	c.Write("\n")
}
