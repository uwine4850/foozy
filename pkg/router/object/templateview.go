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

type JsonObjectTemplateView struct {
	View    IView
	DTO     *rest.DTO
	Message irest.IMessage
}

func (v *JsonObjectTemplateView) Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	if v.View == nil {
		panic("the View field must not be nil")
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

	var filledMessage any
	if v.Message != nil {
		// Retrieves objects by their names and adds them to the general _context map.
		objectContext, err := contextByNameToObjectContext(objectContext[v.View.ObjectsName()[0]])
		if err != nil {
			return func() { v.View.OnError(w, r, manager, err) }
		}
		fmap.MergeMap((*map[string]interface{})(&_context), objectContext)
		manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, _context)
		_filledMessage, err := fillMessage(v.DTO, &_context, v.Message)
		if err != nil {
			return func() { v.View.OnError(w, r, manager, err) }
		}
		filledMessage = _filledMessage
	} else {
		fmap.MergeMap((*map[string]interface{})(&_context), objectContext)
		manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, _context)
		filledMessage = _context
	}
	router.SendJson(filledMessage, w)
	return func() {}
}

type JsonMultipleObjectTemplateView struct {
	View    IView
	DTO     *rest.DTO
	Message irest.IMessage
}

func (v *JsonMultipleObjectTemplateView) Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
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

	contextSliceMap := []ObjectContext{}
	var filledMessages []any
	if v.Message != nil {
		// Retrieves objects by their names and adds them to the general _context map.
		for i := 0; i < len(v.View.ObjectsName()); i++ {
			contextObject, err := contextByNameToObjectContext(objectContext[v.View.ObjectsName()[i]])
			if err != nil {
				return func() { v.View.OnError(w, r, manager, err) }
			}
			// The contextBuff variable is needed so that the data from _context is assigned separately to each object.
			// You cannot copy directly to _context, since this data must be static for each object.
			contextBuff := ObjectContext{}
			fmap.MergeMap((*map[string]interface{})(&contextBuff), _context)
			fmap.MergeMap((*map[string]interface{})(&contextBuff), contextObject)
			contextSliceMap = append(contextSliceMap, contextBuff)
		}
		manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, contextSliceMap)
		for i := 0; i < len(contextSliceMap); i++ {
			filledMessage, err := fillMessage(v.DTO, &contextSliceMap[i], v.Message)
			if err != nil {
				return func() { v.View.OnError(w, r, manager, err) }
			}
			filledMessages = append(filledMessages, filledMessage)
		}
		return func() { router.SendJson(filledMessages, w) }
	} else {
		contextBuff := _context
		fmap.MergeMap((*map[string]interface{})(&contextBuff), objectContext)
		manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, contextBuff)
		contextSliceMap = append(contextSliceMap, contextBuff)
	}
	router.SendJson(contextSliceMap[0], w)
	return func() {}
}

type JsonAllTemplateView struct {
	View    IView
	DTO     *rest.DTO
	Message irest.IMessage
}

func (v *JsonAllTemplateView) Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
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

	contextSliceMap := []ObjectContext{}
	var filledMessages []any
	if v.Message != nil {
		// Retrieves objects by their names and adds them to the general _context map.
		objectContextData := objectContext[v.View.ObjectsName()[0]]
		objectBytes, err := json.Marshal(objectContextData)
		if err != nil {
			return func() { v.View.OnError(w, r, manager, err) }
		}
		var objectContextMap []ObjectContext
		if err := json.Unmarshal(objectBytes, &objectContextMap); err != nil {
			return func() { v.View.OnError(w, r, manager, err) }
		}
		// One object has multiple values.
		// The contextBuff variable is needed so that the data from _context is assigned separately to each object.
		// You cannot copy directly to _context, since this data must be static for each object.
		for i := 0; i < len(objectContextMap); i++ {
			contextBuff := ObjectContext{}
			fmap.MergeMap((*map[string]interface{})(&contextBuff), objectContextMap[i])
			fmap.MergeMap((*map[string]interface{})(&contextBuff), _context)
			contextSliceMap = append(contextSliceMap, contextBuff)
		}
		manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, contextSliceMap[0])
		for i := 0; i < len(contextSliceMap); i++ {
			filledMessage, err := fillMessage(v.DTO, &contextSliceMap[i], v.Message)
			if err != nil {
				return func() { v.View.OnError(w, r, manager, err) }
			}
			filledMessages = append(filledMessages, filledMessage)
		}
		return func() { router.SendJson(filledMessages, w) }
	} else {
		fmap.MergeMap((*map[string]interface{})(&_context), objectContext)
		manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, _context)
		contextSliceMap = append(contextSliceMap, _context)
	}
	router.SendJson(contextSliceMap[0], w)
	return func() {}
}

func contextByNameToObjectContext(contextData interface{}) (ObjectContext, error) {
	objectBytes, err := json.Marshal(contextData)
	if err != nil {
		return nil, err
	}
	var objectContext ObjectContext
	if err := json.Unmarshal(objectBytes, &objectContext); err != nil {
		return nil, err
	}
	return objectContext, nil
}

func fillMessage(dto *rest.DTO, objectContext *ObjectContext, messageType irest.IMessage) (irest.IMessage, error) {
	newMessage := reflect.New(reflect.TypeOf(messageType)).Elem()
	if err := restmapper.FillMessageFromMap((*map[string]interface{})(objectContext), typeopr.Ptr{}.New(&newMessage)); err != nil {
		return nil, err
	}
	newMessageInface := newMessage.Interface().(irest.IMessage)
	if err := rest.DeepCheckSafeMessage(dto, typeopr.Ptr{}.New(&newMessageInface)); err != nil {
		return nil, err
	}
	return newMessageInface, nil
}
