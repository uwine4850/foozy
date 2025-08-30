## object
This package is designed to display data more conveniently on the page using view.<br>
More detailed use of this package is shown in the [tests](https://github.com/uwine4850/foozy/tree/master/tests/router_test/object).

### IView
The interface implements the basic structure of any IView. `ITemplateView` is used to display HTML page in a simpler and more convenient way.<br>
For the view to work correctly, you need to create a new structure (for example MyObjView), embed a ready-made implementation of the view 
(for example ObjView) into it, then you need to initialize this structure in the ITemplateView field in the TemplateView data type.
```golang
type IView interface {
	// Object receives data from the selected table and writes it to a variable structure.
	// IMPORTANT: connects to the database in this method (or others), but closes the connection only in the TemplateView.
	Object(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) (Context, error)
	Context(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) (Context, error)
	Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) (bool, func())
	OnError(w http.ResponseWriter, r *http.Request, manager interfaces.Manager, err error)
	ObjectsName() []string
}
```

### IViewDatabase
Interface that provides object unified access to the database.
Since each object queries the database, it is necessary to unify access 
to the database so as not to be dependent on a particular database.
```golang
type IViewDatabase interface {
	// SelectAll selects all data from the table.
	SelectAll(tableName string) ([]map[string]interface{}, error)
	// SelectWhereEqual selects all data from the table according to the specified condition.
	SelectWhereEqual(tableName string, colName string, val any) ([]map[string]interface{}, error)
}
```

#### GetContext
Retrieves the `Context` from the manager.
It is important to understand that this method can only be used when the IView.Object method has completed running, 
for example in `IView.Context`.
```golang
func GetContext(manager interfaces.Manager) (Context, error) {
	objectInterface, ok := manager.OneTimeData().GetUserContext(namelib.OBJECT.OBJECT_CONTEXT)
	if !ok {
		return nil, errors.New("unable to get object context")
	}
	object := objectInterface.(Context)
	return object, nil
}
```

### AllView
Displays HTML page by passing all data from the selected table to it.
If the `slug` parameter is set, all data from the table that match the condition will be output.<br>
If the `slug` parameter is not set, all data from the table will be output.
```golang
type AllView struct {
	BaseView
	Name       string        `notdef:"true"`
	TableName  string        `notdef:"true"`
	Database   IViewDatabase `notdef:"true"`
	Slug       string
	FillStruct interface{} `notdef:"true"`
}
```

### FormView
This object is designed to process forms sent via HTML forms.
```golang
type FormView struct {
	BaseView

	FormStruct   interface{} `notdef:"true"`
	ValidateCSRF bool
}
```

#### FormInterface
Retrieves the form interface itself from the interface pointer.
```golang
func (v *FormView) FormInterface(manager interfaces.ManagerOneTimeData) (interface{}, error) {
	context, ok := manager.GetUserContext(namelib.OBJECT.OBJECT_CONTEXT)
	if !ok {
		return nil, errors.New("the ObjectContext not found")
	}
	objectContext, ok := context.(Context)
	if !ok {
		return nil, errors.New("the ObjectContext type assertion error")
	}
	return reflect.Indirect(reflect.ValueOf(objectContext[namelib.OBJECT.OBJECT_CONTEXT_FORM])).Interface(), nil
}
```

### MultipleObjectView
Used to display data from multiple sources at once. You also need to use `MultipleObject` to transfer data about a specific object.
```golang
type MultipleObjectView struct {
	BaseView

	Database        IViewDatabase    `notdef:"true"`
	MultipleObjects []MultipleObject `notdef:"true"`
}
```

### ObjView
Displays only the HTML page only with a specific row from the database.<br>
Needs to be used with slug parameter URL path, specify the name of the parameter in the Slug parameter.
```golang
type ObjView struct {
	BaseView

	Name       string        `notdef:"true"`
	TableName  string        `notdef:"true"`
	Database   IViewDatabase `notdef:"true"`
	FillStruct interface{}   `notdef:"true"`
	Slug       string        `notdef:"true"`
}
```

### TemplateView
Renders an HTML page using a template engine and `View`.
```golang
type TemplateView struct {
	TemplatePath string
	View         IView
	isSkipRender bool
}
```

### TemplateRedirectView
Renders an HTML page using a template engine and `View`.<br>
Redirects the page to the selected address.
```golang
type TemplateRedirectView struct {
	View        IView
	RedirectUrl string
}
```

### JsonObjectTemplateView
Function used to display ObjectView as JSON data.<br>
If the Messages field is empty, it renders JSON as a regular `TemplateView`.
```golang
type JsonObjectTemplateView struct {
	View            IView
	DTO             *rest.DTO
	Message         irest.Message
	onMessageFilled OnMessageFilled
}
```

### JsonMultipleObjectTemplateView
Function used to display MultipleObjectView as JSON data.<br>
If the Messages field is empty, it renders JSON as a regular `TemplateView`.
```golang
type JsonMultipleObjectTemplateView struct {
	View            IView
	DTO             *rest.DTO
	Messages        map[string]irest.Message
	onMessageFilled OnMessageFilled
}
```

### JsonAllTemplateView
Function used to display AllView as JSON data.<br>
If the Messages field is empty, it renders JSON as a regular `TemplateView`.
```golang
type JsonAllTemplateView struct {
	View            IView
	DTO             *rest.DTO
	Message         irest.Message
	onMessageFilled OnMessageFilled
}
```