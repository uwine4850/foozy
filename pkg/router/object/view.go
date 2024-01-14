package object

import (
	"github.com/uwine4850/foozy/pkg/interfaces"
	"net/http"
	"reflect"
	"strings"
)

type View interface {
	Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (map[string]interface{}, error)
	Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) map[string]interface{}
	Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func()
	OnError(e func(err error))
}

type UserView interface {
	Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) map[string]interface{}
}

type TemplateStruct struct {
	s reflect.Value
	m map[string]string
}

func (t TemplateStruct) F(name string) string {
	if t.s.Kind() != reflect.Invalid {
		return t.s.FieldByName(strings.ToUpper(string(name[0])) + name[1:]).String()
	}
	if t.m != nil {
		return t.m[name]
	}
	return ""
}
