package runtime

import (
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

type Object interface {
	Type() JSObjectType
}

type JSObjectType string

const (
	JS_OBJECT_TYPE_OBJECT   = "object"
	JS_OBJECT_TYPE_STRING   = "string"
	JS_OBJECT_TYPE_FUNCTION = "function"
)

type JSObject struct {
	properties map[string]Object
}

func (self *JSObject) Type() JSObjectType { return JS_OBJECT_TYPE_OBJECT }

func (self *JSObject) DefineProperty(prop string, value Object) {
	self.properties[prop] = value
}

func (self *JSObject) GetProperty(prop string) (Object, error) {
	obj := self.properties[prop]
	if obj == nil {
		return nil, &ReferenceError{prop}
	}

	return obj, nil
}

type JSString string

func (self JSString) Type() JSObjectType { return JS_OBJECT_TYPE_STRING }

type JSFunction struct {
	fn func([]Object)
}

func (self *JSFunction) FuncName() string {
	fullName := runtime.FuncForPC(reflect.ValueOf(self.fn).Pointer()).Name()
	return strings.TrimPrefix(filepath.Ext(fullName), ".")
}

func (self *JSFunction) Type() JSObjectType { return JS_OBJECT_TYPE_FUNCTION }
