package object

import (
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/utils"
	"net/http"
)

type TemplateView struct {
	TemplatePath string
	View

	onError func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error)
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

func (v *TemplateView) OnError(e func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error)) {
	v.onError = e
}

func (v *TemplateView) Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	if v.View == nil {
		panic("the View field must not be nil")
	}
	permissions, f := v.View.Permissions(w, r, manager)
	if !permissions {
		return func() { f() }
	}
	context := v.View.Context(w, r, manager)
	object, err := v.View.Object(w, r, manager)
	if err != nil {
		return func() { v.onError(w, r, manager, err) }
	}
	utils.MergeMap(&context, object)
	manager.SetContext(context)
	manager.SetTemplatePath(v.TemplatePath)
	err = manager.RenderTemplate(w, r)
	if err != nil {
		return func() { v.onError(w, r, manager, err) }
	}
	return func() {}
}
