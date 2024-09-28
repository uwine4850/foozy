## Object package
This package contains features and interfaces for more convenient template management.
To start working with an object, you need to create a new structure and embed the selected object into it. It is important to note that it 
can be built-in only one object.

### type IView interface
The interface that each View must implement.

* _Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (ObjectContext, error)_ - method 
accesses the database and writes them to the templating context.<br>
* _Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) ObjectContext_ - method that
needs to be overridden in the user structure. The important point is that the __Object__ method writes data to the context
before executing this method, so you need to use the method to get the data that is set in Object 
manager.OneTimeData().GetUserContext(namelib.OBJECT_CONTEXT). An active database connection is also available in this method,
it can be obtained using the GetDB() method.<br>
* _Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func()_ - to use the method, 
it must be overridden. Using this method, you can define access rights for an address. In case of access blocking you 
need to return false and the function to be executed.<br>
* _OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error)_ - the method needs to be overridden.
This method will be executed when an error occurs during the internal execution of the algorithms.<br>
* _ObjectsName() []string_ - returns the names of objects, or one object.

__GetObjectContext__
```
GetObjectContext(manager interfaces.IManager) (ObjectContext, error)
```
Returns from the ObjectContext manager.
It is important to understand that this method can only be used when the IView.Object method has completed its work, 
for example, IView.Context.

## Display View as HTML.

### type TemplateView struct
A structure to launch any IView interface.

* _Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func()_ - a method that starts the execution 
of the http request handler.

## Відображення View як JSON.
Each of these structures has a Call method that is used to run it.
Also, each structure displays data only in JSON format.
If you do not pass data to the `message` field, the JSON will be displayed without DTO validation and in `TemplateView` format.

__type JsonObjectTemplateView struct__<br>
A structure that displays an ObjectView.

__type JsonMultipleObjectTemplateView struct__<br>
A structure that displays an ObjectView. MultipleObjectView.

__type JsonAllTemplateView struct__<br>
A structure that displays an ObjectView. AllView.

## Display objects from the database in various formats.

### type ObjView struct
Detailed view of a specific record from the database. Example of use:
```
type ProfileView struct {
    object.ObjView
}

func (v *ProfileView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (ObjectContext, error) {
    fmt.Println(v.GetContext())
    return ObjectContext{"id": 50000}, nil
}

func Init() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
    db := database.NewDatabase("root", "1111", "localhost", "3406", "foozy")
    view := object.TemplateView{
        TemplatePath: "project/templates/profile.html",
        View: &ProfileView{
            object.ObjView{
                Name:       "profile",
                DB:         db,
                TableName:  "auth",
                FillStruct: User{},
                Slug:       "id",
            },
        },
    }
    return view.Call
}
```
_Name_ - the name by which you can refer to the launch from the database.<br>
_DB_ - instance of the database.<br>
_TableName_ - table name<br>
_FillStruct_ - a structure that describes a table.<br>
_Slug_ - the slug value by which you want to find the value in the table.<br>

### type MultipleObjectView struct
It performs the same actions as __ObjView__ with one difference - finding records in the database occurs in several
tables in the database. Example of use:
```
type ProfileMulView struct {
    object.MultipleObjectView
}

func (v *ProfileMulView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (ObjectContext, error) {
    return ObjectContext{}, nil
}

func Init() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
    db := database.NewDatabase("root", "1111", "localhost", "3406", "foozy")
    view := object.TemplateView{
        TemplatePath: "project/templates/profile.html",
        View: &ProfileMulView{
            object.MultipleObjectView{
                DB: db,
                MultipleObjects: []object.MultipleObject{
                    {
                        Name:       "profile",
                        TaleName:   "auth",
                        SlugName:   "id",
                        SlugField:  "id",
                        FillStruct: User{},
                    },
                    {
                        Name:       "tee",
                        TaleName:   "tee",
                        SlugName:   "tee",
                        SlugField:    "id",
                        FillStruct: Tee{},
                    },
                },
            },
        },
    }
    return view.Call
}
```
_DB_ - instance of the database.<br>
_MultipleObjects_ - instance of the `MultipleObjects` structure.<br>

`MultipleObject`<br>
A structure that represents data about a specific record in the database.

_Name_ - the name by which you can refer to the launch from the database.<br>
_TaleName_ - table name<br>
_SlugName_ - the name of the slug to get its value.<br>
_SlugField_ - the name of the column by which you want to search for values in the table using slug.
_FillStruct_ - a structure that describes a table.<br>

### type AllView struct
Displays all data from the table. Example of use:
```
type ProjectView struct {
    object.AllView
}

func (v *ProjectView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (ObjectContext, error) {
    return ObjectContext{}, nil
}

func Init() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
    db := database.NewDatabase("root", "1111", "localhost", "3406", "foozy")
    view := object.TemplateView{
        TemplatePath: "project/templates/profile.html",
        View: &ProjectView{
            object.AllView{
                Name:       "data",
                DB:         db,
                TableName:  "project",
                FillStruct: Project{},
            },
        },
    }
    return view.Call
}
```

### type FormView struct
This object is slightly different from other objects. The difference is that the main purpose of this 
object - read and process form data. The form is processed in the __Object__ method and then passed to the context.
You can get the value in the `Context` method as usual using `manager.OneTimeData().GetUserContext(namelib.OBJECT_CONTEXT)`.
Also, if necessary, you can display the HTML page, but this is not necessary. In order not to display the page, you need to call it 
`object.TemplateView.SkipRender()` method.
To redirect to another page, you can call `http.Redirect` directly in the `Context` method.

It is important to note that the `NotNilFormFields` parameter is universal. If you pass the sign "*" to it as the first element, all fields
structures will be checked for emptiness. If you specify more structure fields after this sign, they will be excluded
from checking for emptiness. You can also simply pass the fields to be checked without the "*" sign.

If the `ValidateCSRF` parameter is `true`, the CSRF token will be checked.

__FormInterface__
```
FormInterface(manager interfaces.IManagerOneTimeData) (interface{}, error)
```
Returns the completed form in interface format.

Example of use:
```
type ObjectForm struct {
	Text []string        `form:"text"`
	File []form.FormFile `form:"file"`
}

type MyFormView struct {
	object.FormView
}

func (v *MyFormView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (object.ObjectContext, error) {
	filledFormInterface, err := v.FormInterface(manager.OneTimeData())
	if err != nil {
		return nil, err
	}
	filledForm := filledFormInterface.(ObjectForm)
	if filledForm.Text[0] != "field" {
		return nil, errors.New("FormView unexpected text field value")
	}
	if filledForm.File[0].Header.Filename != "x.png" {
		return nil, errors.New("FormView unexpected file field value")
	}
	return object.ObjectContext{}, nil
}

func (v *MyFormView) OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error) {
	panic(err)
}

func MyFormViewHNDL() func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func() {
	tv := object.TemplateView{
		TemplatePath: "",
		View: &MyFormView{
			object.FormView{
				FormStruct:       ObjectForm{},
				NotNilFormFields: []string{"Text", "File"},
				NilIfNotExist:    []string{},
			},
		},
	}
	tv.SkipRender()
	return tv.Call
}
```