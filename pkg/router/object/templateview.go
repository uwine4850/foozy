package object

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/uwine4850/foozy/pkg/debug"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/interfaces/irest"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/rest"
	"github.com/uwine4850/foozy/pkg/router/rest/restmapper"
	"github.com/uwine4850/foozy/pkg/typeopr"
	"github.com/uwine4850/foozy/pkg/utils/fmap"
	"github.com/uwine4850/foozy/pkg/utils/fstruct"
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
	debug.RequestLogginIfEnable(debug.P_OBJECT, "run template view")
	if v.View == nil {
		panic("the ITemplateView field must not be nil")
	}
	defer func() {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "close template view database")
		err := v.View.CloseDb()
		if err != nil {
			debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			v.View.OnError(w, r, manager, err)
		}
	}()
	debug.RequestLogginIfEnable(debug.P_OBJECT, "handle object")
	objectContext, err := v.View.Object(w, r, manager)
	if err != nil {
		return func() {
			debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			v.View.OnError(w, r, manager, err)
		}
	}
	debug.RequestLogginIfEnable(debug.P_OBJECT, "handle context")
	manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, objectContext)
	_context, err := v.View.Context(w, r, manager)
	if err != nil {
		return func() {
			debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			v.View.OnError(w, r, manager, err)
		}
	}
	fmap.MergeMap((*map[string]interface{})(&objectContext), _context)
	manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, objectContext)

	if v.isSkipRender {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "skip render")
		return func() {}
	}

	debug.RequestLogginIfEnable(debug.P_OBJECT, "handle permissions")
	permissions, f := v.View.Permissions(w, r, manager)
	if !permissions {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "permissions are not granted")
		return func() { f() }
	}
	manager.Render().SetContext(objectContext)
	manager.Render().SetTemplatePath(v.TemplatePath)
	err = manager.Render().RenderTemplate(w, r)
	if err != nil {
		return func() {
			debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			v.View.OnError(w, r, manager, err)
		}
	}
	return func() {}
}

// TemplateRedirectView processes the object.
// Redirects the page to the selected address.
type TemplateRedirectView struct {
	View        IView
	RedirectUrl string
}

func (v *TemplateRedirectView) Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	debug.RequestLogginIfEnable(debug.P_OBJECT, "run template view")
	if v.View == nil {
		panic("the ITemplateView field must not be nil")
	}
	defer func() {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "close template view database")
		err := v.View.CloseDb()
		if err != nil {
			debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			v.View.OnError(w, r, manager, err)
		}
	}()
	debug.RequestLogginIfEnable(debug.P_OBJECT, "handle object")
	objectContext, err := v.View.Object(w, r, manager)
	if err != nil {
		return func() {
			debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			v.View.OnError(w, r, manager, err)
		}
	}
	debug.RequestLogginIfEnable(debug.P_OBJECT, "handle context")
	manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, objectContext)
	_context, err := v.View.Context(w, r, manager)
	if err != nil {
		return func() {
			debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			v.View.OnError(w, r, manager, err)
		}
	}
	fmap.MergeMap((*map[string]interface{})(&objectContext), _context)
	manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, objectContext)

	debug.RequestLogginIfEnable(debug.P_OBJECT, "handle permissions")
	permissions, f := v.View.Permissions(w, r, manager)
	if !permissions {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "permissions are not granted")
		return func() { f() }
	}
	if v.RedirectUrl == "" {
		return func() { http.Redirect(w, r, r.URL.Path, http.StatusFound) }
	}
	return func() { http.Redirect(w, r, v.RedirectUrl, http.StatusFound) }
}

// JsonObjectTemplateView is used to display ObjectView as JSON data.
// If the Messages field is empty, it renders JSON as a regular TemplateView.
type JsonObjectTemplateView struct {
	View            IView
	DTO             *rest.DTO
	Message         irest.IMessage
	onMessageFilled func(message any, manager interfaces.IManager) error
}

func (v *JsonObjectTemplateView) Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	defer func() {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "close view database")
		err := v.View.CloseDb()
		if err != nil {
			v.View.OnError(w, r, manager, err)
			debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			return
		}
	}()
	debug.RequestLogginIfEnable(debug.P_OBJECT, "run JsonObjectTemplateView")
	onError, viewObject, viewContext := baseParseView(v.View, w, r, manager)
	if onError != nil {
		return onError
	}
	var filledMessage any
	if v.Message != nil {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "fill DTO message...")
		// Retrieves objects by their names and adds them to the general viewContext map.
		objectContext, err := contextByNameToObjectContext(viewObject[v.View.ObjectsName()[0]])
		if err != nil {
			return func() {
				debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
				v.View.OnError(w, r, manager, err)
			}
		}
		fmap.MergeMap((*map[string]interface{})(&viewContext), objectContext)
		manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, viewContext)
		_filledMessage, err := fillMessage(v.DTO, &viewContext, v.Message)
		if err != nil {
			return func() {
				debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
				v.View.OnError(w, r, manager, err)
			}
		}
		if v.onMessageFilled != nil {
			messageType := reflect.TypeOf(v.Message)
			tempMessage := reflect.New(messageType)
			tempMessage.Elem().Set(reflect.ValueOf(_filledMessage))
			if err := v.onMessageFilled(tempMessage.Interface(), manager); err != nil {
				return func() {
					debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
					v.View.OnError(w, r, manager, err)
				}
			}
			_filledMessage = tempMessage.Elem().Interface().(irest.IMessage)
		}
		filledMessage = _filledMessage
	} else {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "pass context without DTO message")
		fmap.MergeMap((*map[string]interface{})(&viewContext), viewObject)
		manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, viewContext)
		filledMessage = viewContext
	}
	debug.RequestLogginIfEnable(debug.P_OBJECT, "send json")
	router.SendJson(filledMessage, w)
	return func() {}
}

func (v *JsonObjectTemplateView) OnMessageFilled(fn func(message any, manager interfaces.IManager) error) {
	v.onMessageFilled = fn
}

// JsonMultipleObjectTemplateView is used to display MultipleObjectView as JSON data.
// If the Messages field is empty, it renders JSON as a regular TemplateView.
type JsonMultipleObjectTemplateView struct {
	View    IView
	DTO     *rest.DTO
	Message irest.IMessage
}

func (v *JsonMultipleObjectTemplateView) Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	defer func() {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "close view database")
		err := v.View.CloseDb()
		if err != nil {
			v.View.OnError(w, r, manager, err)
			debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			return
		}
	}()
	debug.RequestLogginIfEnable(debug.P_OBJECT, "run JsonMultipleObjectTemplateView")
	onError, viewObject, viewContext := baseParseView(v.View, w, r, manager)
	if onError != nil {
		return onError
	}

	contextSliceMap := []Context{}
	var filledMessages []any
	if v.Message != nil {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "fill DTO messages")
		// Retrieves objects by their names and adds them to the general viewContext map.
		for i := 0; i < len(v.View.ObjectsName()); i++ {
			viewObjectContext, err := contextByNameToObjectContext(viewObject[v.View.ObjectsName()[i]])
			if err != nil {
				return func() {
					debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
					v.View.OnError(w, r, manager, err)
				}
			}
			// The contextBuff variable is needed so that the data from viewContext is assigned separately to each object.
			// You cannot copy directly to viewContext, since this data must be static for each object.
			contextBuff := Context{}
			fmap.MergeMap((*map[string]interface{})(&contextBuff), viewContext)
			fmap.MergeMap((*map[string]interface{})(&contextBuff), viewObjectContext)
			contextSliceMap = append(contextSliceMap, contextBuff)
		}
		manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, contextSliceMap)
		for i := 0; i < len(contextSliceMap); i++ {
			filledMessage, err := fillMessage(v.DTO, &contextSliceMap[i], v.Message)
			if err != nil {
				return func() {
					debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
					v.View.OnError(w, r, manager, err)
				}
			}
			filledMessages = append(filledMessages, filledMessage)
		}
		return func() { router.SendJson(filledMessages, w) }
	} else {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "pass context without DTO messages")
		fmap.MergeMap((*map[string]interface{})(&viewContext), viewObject)
		manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, viewContext)
		contextSliceMap = append(contextSliceMap, viewContext)
	}
	debug.RequestLogginIfEnable(debug.P_OBJECT, "send json")
	router.SendJson(contextSliceMap[0], w)
	return func() {}
}

// JsonAllTemplateView is used to display AllView as JSON data.
// If the Messages field is empty, it renders JSON as a regular TemplateView.
type JsonAllTemplateView struct {
	View    IView
	DTO     *rest.DTO
	Message irest.IMessage
}

func (v *JsonAllTemplateView) Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	defer func() {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "close view database")
		err := v.View.CloseDb()
		if err != nil {
			v.View.OnError(w, r, manager, err)
			debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			return
		}
	}()
	debug.RequestLogginIfEnable(debug.P_OBJECT, "run JsonAllTemplateView")
	onError, viewObject, viewContext := baseParseView(v.View, w, r, manager)
	if onError != nil {
		return onError
	}

	contextSliceMap := []Context{}
	var filledMessages []any
	if v.Message != nil {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "fill DTO messages")
		// Retrieves objects by their names and adds them to the general viewContext map.
		objectBytes, err := json.Marshal(viewObject[v.View.ObjectsName()[0]])
		if err != nil {
			return func() {
				debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
				v.View.OnError(w, r, manager, err)
			}
		}
		var objectContextMap []Context
		if err := json.Unmarshal(objectBytes, &objectContextMap); err != nil {
			return func() {
				debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
				v.View.OnError(w, r, manager, err)
			}
		}
		// One object has multiple values.
		// The contextBuff variable is needed so that the data from viewContext is assigned separately to each object.
		// You cannot copy directly to viewContext, since this data must be static for each object.
		for i := 0; i < len(objectContextMap); i++ {
			contextBuff := Context{}
			fmap.MergeMap((*map[string]interface{})(&contextBuff), objectContextMap[i])
			fmap.MergeMap((*map[string]interface{})(&contextBuff), viewContext)
			contextSliceMap = append(contextSliceMap, contextBuff)
		}
		manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, contextSliceMap[0])
		for i := 0; i < len(contextSliceMap); i++ {
			filledMessage, err := fillMessage(v.DTO, &contextSliceMap[i], v.Message)
			if err != nil {
				return func() {
					debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
					v.View.OnError(w, r, manager, err)
				}
			}
			filledMessages = append(filledMessages, filledMessage)
		}
		return func() { router.SendJson(filledMessages, w) }
	} else {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "pass context without DTO messages")
		fmap.MergeMap((*map[string]interface{})(&viewContext), viewObject)
		manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, viewContext)
		contextSliceMap = append(contextSliceMap, viewContext)
	}
	debug.RequestLogginIfEnable(debug.P_OBJECT, "send json")
	router.SendJson(contextSliceMap[0], w)
	return func() {}
}

func baseParseView(view IView, w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (onError func(), viewObject Context, viewContext Context) {
	if view == nil {
		panic("the ITemplateView field must not be nil")
	}
	realView := reflect.ValueOf(getRealView(view))
	if err := fstruct.CheckNotDefaultFields(typeopr.Ptr{}.New(&realView)); err != nil {
		onError = func() {
			debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			view.OnError(w, r, manager, err)
		}
		return
	}
	var err error
	debug.RequestLogginIfEnable(debug.P_OBJECT, "handle object")
	viewObject, err = view.Object(w, r, manager)
	if err != nil {
		onError = func() {
			debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			view.OnError(w, r, manager, err)
		}
		return
	}
	manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, viewObject)

	debug.RequestLogginIfEnable(debug.P_OBJECT, "handle context")
	viewContext, err = view.Context(w, r, manager)
	if err != nil {
		onError = func() {
			debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			view.OnError(w, r, manager, err)
		}
		return
	}

	debug.RequestLogginIfEnable(debug.P_OBJECT, "handle permissions")
	permissions, f := view.Permissions(w, r, manager)
	if !permissions {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "permissions are not granted")
		onError = func() { f() }
		return
	}
	return
}

func getRealView(wrapperView IView) IView {
	rViewValue := reflect.ValueOf(wrapperView).Elem()
	rViewType := reflect.TypeOf(wrapperView).Elem()
	for i := 0; i < rViewType.NumField(); i++ {
		fieldValue := rViewValue.Field(i).Addr()
		if typeopr.IsImplementInterface(typeopr.Ptr{}.New(&fieldValue), (*IView)(nil)) {
			return fieldValue.Interface().(IView)
		}
	}
	return nil
}

// contextByNameToObjectContext converts the View context data into an ObjectContext object.
func contextByNameToObjectContext(contextData interface{}) (Context, error) {
	objectBytes, err := json.Marshal(contextData)
	if err != nil {
		return nil, err
	}
	var objectContext Context
	if err := json.Unmarshal(objectBytes, &objectContext); err != nil {
		return nil, err
	}
	return objectContext, nil
}

// fillMessage fills the DTO Message with the data passed to the objectContext.
// It is important to highlight that the messageType argument is used only to obtain the message type; an instance of this type is returned.
func fillMessage(dto *rest.DTO, objectContext *Context, messageType irest.IMessage) (irest.IMessage, error) {
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
