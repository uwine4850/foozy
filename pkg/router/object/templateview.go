package object

import (
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/utils"
	"net/http"
)

type TemplateView struct {
	UserView
	TemplatePath string
	View

	onError func(err error)
}

func (v *TemplateView) Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func()) {
	return true, func() {}
}

func (v *TemplateView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) map[string]interface{} {
	return map[string]interface{}{}
}

// Object sets a slice of rows from the database.
func (v *TemplateView) Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (v *TemplateView) OnError(e func(err error)) {
	v.onError = e
}

func (v *TemplateView) Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	if v.UserView == nil {
		panic("the UserView field must not be nil")
	}
	permissions, f := v.UserView.Permissions(w, r, manager)
	if !permissions {
		return func() { f() }
	}
	context := v.UserView.Context(w, r, manager)
	object, err := v.View.Object(w, r, manager)
	if err != nil {
		return func() { v.onError(err) }
	}
	utils.MergeMap(&context, object)
	manager.SetContext(context)
	manager.SetTemplatePath(v.TemplatePath)
	err = manager.RenderTemplate(w, r)
	if err != nil {
		return func() { v.onError(err) }
	}
	return func() {}
}
