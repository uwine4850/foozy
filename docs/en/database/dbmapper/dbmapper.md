## package dbmapper
Writes data to the selected object. Can perform various operations on 
filled data.

You can see how the package works in these [tests](https://github.com/uwine4850/foozy/tree/master/tests/database/dbmapper_test).

### type Mapper struct
A structure used to populate an object with data from the database.
A list of data from the database is passed to `DatabaseResult`. 
In the `Output` field you need pass a reference to the desired object.

It is important to note that the data must be in the form of a slice.

At the moment, three types of data can be used in the `Output` field:
* Structure.
* The structure in the `reflect.Value` statement.
* Data type map[string]string.

If a structure is used as `Output`, it must be properly configured.
Example structure:
```
type DbTestMapper struct {
	Col1 string `db:"col1"`
	Col2 string `db:"col2"`
	Col3 string `db:"col3" empty:"-err"`
	Col4 string `db:"col4" empty:"0"`
}
```
Two tags can be used for setting:
* db — is a required tag. Needed to know which field corresponds to which
column in the table. Accordingly, it must specify the name of the column.
* empty — optional tag. This tag is used for the what command 
must be done when the field is empty. There are currently only two values:
    * -err — if the field is empty, outputs the corresponding error.
    * plain text — text that will replace empty values.

__Fill__
```
Fill() error 
```
Fills in the data.

### Other features of the package.

__FillStructFromDb__
```
FillStructFromDb(dbRes map[string]interface{}, fillPtr itypeopr.IPtr) error
```
Fills the structure with data from the database.
Each variable of the placeholder structure must have a `db` tag that corresponds to the name of the column in the database, for example `db: "name"`. You can also use the `empty` tag, which is described in more detail above.

__FillMapFromDb__
```
FillMapFromDb(dbRes map[string]interface{}, fill *map[string]string) error
```
Fills the map with data from the database.

__FillReflectValueFromDb__
```
FillReflectValueFromDb(dbRes map[string]interface{}, fill *reflect.Value) error
```
Fills a structure whose type is *reflect.Value. That is, the method fills the data from the database into the structure, which is created with the help 
package reflect.

__ParamsValueFromStruct__
```
ParamsValueFromStruct(filledStructurePtr itypeopr.IPtr, nilIfEmpty []string) (map[string]any, error)
```
Creates a map from a structure that describes a table.
You need a completed structure to work correctly, and required fields must have the `db:"<column name>"` tag.
You can also use the `empty` tag, which is described in more detail above.