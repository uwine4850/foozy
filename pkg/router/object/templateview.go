package object

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"github.com/uwine4850/foozy/pkg/debug"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/interfaces/irest"
	"github.com/uwine4850/foozy/pkg/mapper"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/router/rest"
	"github.com/uwine4850/foozy/pkg/typeopr"
	"github.com/uwine4850/foozy/pkg/utils/fmap"
	"github.com/uwine4850/foozy/pkg/utils/fstruct"
)

type OnMessageFilled func(message any, manager interfaces.IManager) error

type TemplateView struct {
	TemplatePath string
	View         IView
	isSkipRender bool
}

func (v *TemplateView) SkipRender() {
	v.isSkipRender = true
}

func (v *TemplateView) Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
	debug.RequestLogginIfEnable(debug.P_OBJECT, "run template view")
	if v.View == nil {
		panic("the ITemplateView field must not be nil")
	}
	debug.RequestLogginIfEnable(debug.P_OBJECT, "handle object")
	objectContext, err := v.View.Object(w, r, manager)
	if err != nil {
		debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
		v.View.OnError(w, r, manager, err)
		return nil
	}
	debug.RequestLogginIfEnable(debug.P_OBJECT, "handle context")
	manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, objectContext)
	_context, err := v.View.Context(w, r, manager)
	if err != nil {
		debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
		v.View.OnError(w, r, manager, err)
		return nil
	}
	fmap.MergeMap((*map[string]interface{})(&objectContext), _context)
	manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, objectContext)

	if v.isSkipRender {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "skip render")
		return nil
	} else if manager.Render() == nil {
		return errors.New("IRender is used, but it is not set")
	}

	debug.RequestLogginIfEnable(debug.P_OBJECT, "handle permissions")
	permissions, f := v.View.Permissions(w, r, manager)
	if !permissions {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "permissions are not granted")
		f()
		return nil
	}
	manager.Render().SetContext(objectContext)
	manager.Render().SetTemplatePath(v.TemplatePath)
	err = manager.Render().RenderTemplate(w, r)
	if err != nil {
		debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
		v.View.OnError(w, r, manager, err)
		return nil
	}
	return nil
}

// TemplateRedirectView processes the object.
// Redirects the page to the selected address.
type TemplateRedirectView struct {
	View        IView
	RedirectUrl string
}

func (v *TemplateRedirectView) Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
	debug.RequestLogginIfEnable(debug.P_OBJECT, "run template view")
	if v.View == nil {
		panic("the ITemplateView field must not be nil")
	}
	debug.RequestLogginIfEnable(debug.P_OBJECT, "handle object")
	objectContext, err := v.View.Object(w, r, manager)
	if err != nil {
		debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
		v.View.OnError(w, r, manager, err)
		return nil
	}
	debug.RequestLogginIfEnable(debug.P_OBJECT, "handle context")
	manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, objectContext)
	_context, err := v.View.Context(w, r, manager)
	if err != nil {
		debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
		v.View.OnError(w, r, manager, err)
		return nil
	}
	fmap.MergeMap((*map[string]interface{})(&objectContext), _context)
	manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, objectContext)

	debug.RequestLogginIfEnable(debug.P_OBJECT, "handle permissions")
	permissions, f := v.View.Permissions(w, r, manager)
	if !permissions {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "permissions are not granted")
		f()
		return nil
	}
	if v.RedirectUrl == "" {
		http.Redirect(w, r, r.URL.Path, http.StatusFound)
		return nil
	}
	http.Redirect(w, r, v.RedirectUrl, http.StatusFound)
	return nil
}

// JsonObjectTemplateView is used to display ObjectView as JSON data.
// If the Messages field is empty, it renders JSON as a regular TemplateView.
type JsonObjectTemplateView struct {
	View            IView
	DTO             *rest.DTO
	Message         irest.IMessage
	onMessageFilled OnMessageFilled
}

func (v *JsonObjectTemplateView) Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
	debug.RequestLogginIfEnable(debug.P_OBJECT, "run JsonObjectTemplateView")
	viewObject, viewContext, err := baseParseView(v.View, w, r, manager)
	if err != nil {
		v.View.OnError(w, r, manager, err)
		return nil
	}
	var filledMessage any
	if v.Message != nil {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "fill DTO message...")
		// Retrieves objects by their names and adds them to the general viewContext map.
		objectContext, err := contextByNameToObjectContext(viewObject[v.View.ObjectsName()[0]])
		if err != nil {
			debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			v.View.OnError(w, r, manager, err)
			return nil
		}
		fmap.MergeMap((*map[string]interface{})(&viewContext), objectContext)
		manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, viewContext)
		_filledMessage, err := fillMessage(v.DTO, &viewContext, v.Message)
		if err != nil {
			debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			v.View.OnError(w, r, manager, err)
			return nil
		}
		if v.onMessageFilled != nil {
			tempMessage := makePointerToFilledMessage(v.Message, reflect.ValueOf(_filledMessage))
			if err := runOnMessageFilledFunction(v.onMessageFilled, &_filledMessage, tempMessage, manager); err != nil {
				debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
				v.View.OnError(w, r, manager, err)
				return nil
			}
		}
		filledMessage = _filledMessage
	} else {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "pass context without DTO message")
		fmap.MergeMap((*map[string]interface{})(&viewContext), viewObject)
		manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, viewContext)
		filledMessage = viewContext
	}
	debug.RequestLogginIfEnable(debug.P_OBJECT, "send json")
	router.SendJson(filledMessage, w, http.StatusOK)
	return nil
}

func (v *JsonObjectTemplateView) OnMessageFilled(fn func(message any, manager interfaces.IManager) error) {
	v.onMessageFilled = fn
}

// JsonMultipleObjectTemplateView is used to display MultipleObjectView as JSON data.
// If the Messages field is empty, it renders JSON as a regular TemplateView.
type JsonMultipleObjectTemplateView struct {
	View            IView
	DTO             *rest.DTO
	Messages        map[string]irest.IMessage
	onMessageFilled OnMessageFilled
}

func (v *JsonMultipleObjectTemplateView) Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
	debug.RequestLogginIfEnable(debug.P_OBJECT, "run JsonMultipleObjectTemplateView")
	viewObject, viewContext, err := baseParseView(v.View, w, r, manager)
	if err != nil {
		v.View.OnError(w, r, manager, err)
		return nil
	}

	returnData := Context{}

	if v.Messages != nil {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "fill DTO messages")
		objectsData := Context{}
		fmap.MergeMap((*map[string]interface{})(&objectsData), viewContext)
		fmap.MergeMap((*map[string]interface{})(&objectsData), viewObject)
		manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, objectsData)

		// Fill messages.
		for objectName, message := range v.Messages {
			objectData := objectsData[objectName]
			viewObjectContext, err := contextByNameToObjectContext(objectData)
			if err != nil {
				debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
				v.View.OnError(w, r, manager, err)
				return nil
			}
			filledMessage, err := fillMessage(v.DTO, &viewObjectContext, message)
			if err != nil {
				debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
				v.View.OnError(w, r, manager, err)
				return nil
			}
			if v.onMessageFilled != nil {
				tempMessage := makePointerToFilledMessage(message, reflect.ValueOf(filledMessage))
				if err := runOnMessageFilledFunction(v.onMessageFilled, &filledMessage, tempMessage, manager); err != nil {
					debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
					v.View.OnError(w, r, manager, err)
					return nil
				}
			}
			returnData[objectName] = filledMessage
		}
	} else {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "pass context without DTO messages")
		fmap.MergeMap((*map[string]interface{})(&viewContext), viewObject)
		manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, viewContext)
		returnData = viewContext
	}
	debug.RequestLogginIfEnable(debug.P_OBJECT, "send json")
	router.SendJson(returnData, w, http.StatusOK)
	return nil
}

func (v *JsonMultipleObjectTemplateView) OnMessageFilled(fn func(message any, manager interfaces.IManager) error) {
	v.onMessageFilled = fn
}

// JsonAllTemplateView is used to display AllView as JSON data.
// If the Messages field is empty, it renders JSON as a regular TemplateView.
type JsonAllTemplateView struct {
	View            IView
	DTO             *rest.DTO
	Message         irest.IMessage
	onMessageFilled OnMessageFilled
}

func (v *JsonAllTemplateView) Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
	debug.RequestLogginIfEnable(debug.P_OBJECT, "run JsonAllTemplateView")
	viewObject, viewContext, err := baseParseView(v.View, w, r, manager)
	if err != nil {
		v.View.OnError(w, r, manager, err)
		return nil
	}
	contextSliceMap := []Context{}
	var filledMessages []any
	if v.Message != nil {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "fill DTO messages")
		// Retrieves objects by their names and adds them to the general viewContext map.
		objectBytes, err := json.Marshal(viewObject[v.View.ObjectsName()[0]])
		if err != nil {
			debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			v.View.OnError(w, r, manager, err)
			return nil
		}
		var objectContextMap []Context
		if err := json.Unmarshal(objectBytes, &objectContextMap); err != nil {
			debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
			v.View.OnError(w, r, manager, err)
			return nil
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
		manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, contextSliceMap)
		for i := 0; i < len(contextSliceMap); i++ {
			filledMessage, err := fillMessage(v.DTO, &contextSliceMap[i], v.Message)
			if err != nil {
				debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
				v.View.OnError(w, r, manager, err)
				return nil
			}
			if v.onMessageFilled != nil {
				tempMessage := makePointerToFilledMessage(v.Message, reflect.ValueOf(filledMessage))
				if err := runOnMessageFilledFunction(v.onMessageFilled, &filledMessage, tempMessage, manager); err != nil {
					debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
					v.View.OnError(w, r, manager, err)
					return nil
				}
			}
			filledMessages = append(filledMessages, filledMessage)
		}
		router.SendJson(filledMessages, w, http.StatusOK)
		return nil
	} else {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "pass context without DTO messages")
		fmap.MergeMap((*map[string]interface{})(&viewContext), viewObject)
		manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, viewContext)
		contextSliceMap = append(contextSliceMap, viewContext)
	}
	debug.RequestLogginIfEnable(debug.P_OBJECT, "send json")
	router.SendJson(contextSliceMap[0], w, http.StatusOK)
	return nil
}

func (v *JsonAllTemplateView) OnMessageFilled(fn func(message any, manager interfaces.IManager) error) {
	v.onMessageFilled = fn
}

func baseParseView(view IView, w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (viewObject Context, viewContext Context, outErr error) {
	if view == nil {
		panic("the ITemplateView field must not be nil")
	}
	realView := reflect.ValueOf(getRealView(view))
	if err := fstruct.CheckNotDefaultFields(typeopr.Ptr{}.New(&realView)); err != nil {
		debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
		return nil, nil, err
	}
	var err error
	debug.RequestLogginIfEnable(debug.P_OBJECT, "handle object")
	viewObject, err = view.Object(w, r, manager)
	if err != nil {
		debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
		return nil, nil, err
	}
	manager.OneTimeData().SetUserContext(namelib.OBJECT.OBJECT_CONTEXT, viewObject)

	debug.RequestLogginIfEnable(debug.P_OBJECT, "handle context")
	viewContext, err = view.Context(w, r, manager)
	if err != nil {
		debug.RequestLogginIfEnable(debug.P_ERROR, err.Error())
		return nil, nil, err
	}

	debug.RequestLogginIfEnable(debug.P_OBJECT, "handle permissions")
	permissions, f := view.Permissions(w, r, manager)
	if !permissions {
		debug.RequestLogginIfEnable(debug.P_OBJECT, "permissions are not granted")
		f()
		return nil, nil, err
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
	if err := mapper.DeepCheckDTOSafeMessage(dto, typeopr.Ptr{}.New(&messageType)); err != nil {
		return nil, err
	}
	newMessage := reflect.New(reflect.TypeOf(messageType)).Elem()
	if err := mapper.FillDTOMessageFromMap(*(*map[string]interface{})(objectContext), &newMessage); err != nil {
		return nil, err
	}
	newMessageInface := newMessage.Interface().(irest.IMessage)
	return newMessageInface, nil
}

func makePointerToFilledMessage(messageInstance irest.IMessage, filledMessage reflect.Value) interface{} {
	messageType := reflect.TypeOf(messageInstance)
	tempMessage := reflect.New(messageType)
	tempMessage.Elem().Set(filledMessage)
	return tempMessage.Interface()
}

func runOnMessageFilledFunction(onMessageFilledFn OnMessageFilled, sourceFilledMessage any, filledMessagePointer any, manager interfaces.IManager) error {
	if !typeopr.IsPointer(sourceFilledMessage) {
		return typeopr.ErrValueNotPointer{Value: "sourceFilledMessage"}
	}
	if !typeopr.IsPointer(filledMessagePointer) {
		return typeopr.ErrValueNotPointer{Value: "filledMessagePointer"}
	}
	if err := onMessageFilledFn(filledMessagePointer, manager); err != nil {
		return err
	}
	reflect.ValueOf(sourceFilledMessage).Elem().Set(reflect.ValueOf(filledMessagePointer))
	return nil
}
