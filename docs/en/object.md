## Form package
This package contains features and interfaces for more convenient template management.

__type TemplateView struct__

A structure to launch any IView interface.

* _Call(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) func()_ - a method that starts the execution 
of the http request handler.

__type IView interface__

The interface that each View must implement.

* _Object(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (map[string]interface{}, error)_ - method 
accesses the database and writes them to the templating context.<br>
* _GetContext() map[string]interface{}_ - returns the context that will be installed in the templating tool.<br>
* _Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) map[string]interface{}_ - method that
needs to be overridden in the user structure. An important point is that the __Object__ method writes data to the context
before executing this method, so you can use the __GetContext__ method to get an object from the database.<br>
* _Permissions(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) (bool, func()_ - to use the method, 
it must be overridden. Using this method, you can define access rights for an address. In case of access blocking you 
need to return false and the function to be executed.<br>
* _OnError(w http.ResponseWriter, r *http.Request, manager interfaces.IManager, err error)_ - the method needs to be overridden.
This method will be executed when an error occurs during the internal execution of the algorithms.<br>

## type ObjView struct

Detailed view of a specific record from the database. Example of use:
```
type ProfileView struct {
    object.ObjView
}

func (v *ProfileView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) map[string]interface{} {
    fmt.Println(v.GetContext())
    return map[string]interface{}{"id": 50000}
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

## type MultipleObjectView struct

It performs the same actions as __ObjView__ with one difference - finding records in the database occurs in several
tables in the database. Example of use:
```
type ProfileMulView struct {
    object.MultipleObjectView
}

func (v *ProfileMulView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) map[string]interface{} {
    return map[string]interface{}{}
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
                        "profile",
                        "auth",
                        "id",
                        "id",
                        User{},
                    },
                    {
                        Name:       "tee",
                        TaleName:   "tee",
                        SlugName:   "tee",
                        AIField:    "id",
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
_MultipleObjects_ - instance of the __MultipleObjects__ structure.<br>

__MultipleObject__

A structure that represents data about a specific record in the database.

_Name_ - the name by which you can refer to the launch from the database.<br>
_TaleName_ - table name<br>
_SlugName_ - the name of the slug to get its value.<br>
_AIField_ - the name of the column by which you want to search for values in the table using slug.
_FillStruct_ - a structure that describes a table.<br>

## type AllView struct

Displays all data from the table. Example of use:
```
type ProjectView struct {
    object.AllView
}

func (v *ProjectView) Context(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) map[string]interface{} {
    return map[string]interface{}{}
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
