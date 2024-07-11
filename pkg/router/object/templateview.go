package object

import (
	"net/http"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/utils/fmap"
)

type TemplateView struct {
	TemplatePath string
	View         IView
}

func (v *TemplateView) Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	if v.View == nil {
		panic("the ITemplateView field must not be nil")
	}
	defer func() {
		err := v.View.GetDB().Close()
		if err != nil {
			v.View.OnError(w, r, manager, err)
		}
	}()
	objectContext, err := v.View.Object(w, r, manager)
	if err != nil {
		return func() { v.View.OnError(w, r, manager, err) }
	}
	manager.OneTimeData().SetUserContext(namelib.OBJECT_CONTEXT, objectContext)
	_context, err := v.View.Context(w, r, manager)
	if err != nil {
		return func() { v.View.OnError(w, r, manager, err) }
	}
	fmap.MergeMap((*map[string]interface{})(&objectContext), _context)
	manager.OneTimeData().SetUserContext(namelib.OBJECT_CONTEXT, objectContext)
	permissions, f := v.View.Permissions(w, r, manager)
	if !permissions {
		return func() { f() }
	}
	manager.Render().SetContext(objectContext)
	manager.Render().SetTemplatePath(v.TemplatePath)
	err = manager.Render().RenderTemplate(w, r)
	if err != nil {
		return func() { v.View.OnError(w, r, manager, err) }
	}
	return func() {}
}
