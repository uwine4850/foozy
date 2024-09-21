package object

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/interfaces/irest"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/rest"
	"github.com/uwine4850/foozy/pkg/router/rest/restmapper"
	"github.com/uwine4850/foozy/pkg/typeopr"
	"github.com/uwine4850/foozy/pkg/utils/fmap"
)

type TemplateView struct {
	TemplatePath string
	View         IView
	isSkipRender bool
}

func (v *TemplateView) SkipRender() {
	v.isSkipRender = true
}

func (v *TemplateView) Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	if v.View == nil {
		panic("the ITemplateView field must not be nil")
	}
	defer func() {
		err := v.View.CloseDb()
		if err != nil {
			v.View.OnError(w, r, manager, err)
		}
	}()
	objectContext, err := v.View.Object(w, r, manager)
	if err != nil {
		return func() { v.View.OnError(w, r, manager, err) }
	}
	manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, objectContext)
	_context, err := v.View.Context(w, r, manager)
	if err != nil {
		return func() { v.View.OnError(w, r, manager, err) }
	}
	fmap.MergeMap((*map[string]interface{})(&objectContext), _context)
	manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, objectContext)

	if v.isSkipRender {
		return func() {}
	}

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

type JsonTemplateView struct {
	View    IView
	DTO     *rest.DTO
	Message irest.IMessage
}

func (v *JsonTemplateView) Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	if v.View == nil {
		panic("the ITemplateView field must not be nil")
	}
	defer func() {
		err := v.View.CloseDb()
		if err != nil {
			v.View.OnError(w, r, manager, err)
		}
	}()
	objectContext, err := v.View.Object(w, r, manager)
	if err != nil {
		return func() { v.View.OnError(w, r, manager, err) }
	}
	manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, objectContext)
	_context, err := v.View.Context(w, r, manager)
	if err != nil {
		return func() { v.View.OnError(w, r, manager, err) }
	}
	permissions, f := v.View.Permissions(w, r, manager)
	if !permissions {
		return func() { f() }
	}

	// Retrieves objects by their names and adds them to the general _context map.
	for i := 0; i < len(v.View.ObjectsName()); i++ {
		objectContextData := objectContext[v.View.ObjectsName()[i]]
		objectBytes, err := json.Marshal(objectContextData)
		if err != nil {
			return func() { v.View.OnError(w, r, manager, err) }
		}
		var objectContextMap map[string]any
		if err := json.Unmarshal(objectBytes, &objectContextMap); err != nil {
			return func() { v.View.OnError(w, r, manager, err) }
		}
		fmap.MergeMap((*map[string]interface{})(&_context), objectContextMap)
	}
	manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, _context)

	if v.Message != nil {
		newMessage := reflect.New(reflect.TypeOf(v.Message)).Elem()
		if err := restmapper.FillMessageFromMap((*map[string]interface{})(&_context), typeopr.Ptr{}.New(&newMessage)); err != nil {
			return func() { v.View.OnError(w, r, manager, err) }
		}
		newMessageInface := newMessage.Interface().(irest.IMessage)
		if err := rest.DeepCheckSafeMessage(v.DTO, typeopr.Ptr{}.New(&newMessageInface)); err != nil {
			return func() { v.View.OnError(w, r, manager, err) }
		}
		return func() { router.SendJson(newMessageInface, w) }
	}
	router.SendJson(_context, w)
	return func() {}
}
