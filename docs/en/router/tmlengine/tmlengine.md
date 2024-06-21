## Package tmlengine
This package implements a template generator based on the ``pongo2`` library.<br>
The package can be divided into two modules: logic of the templater and logic of filters for it.

## Template engine
The template engine is used for convenient output of HTML code on a page in a browser. You can read more about pongo2 
in the official repository at the following link https://github.com/flosch/pongo2.

__SetPath__
```
SetPath(path string)
```
Mandatory method. Sets the path to the HTML template.

__SetContext__
```
SetContext(data map[string]interface{})
```
Setting the templating context. Access to variables in HTML is carried out in the following way ``{{ name }}``. If the variables <name>, 
for example, the structure - access to variables will be ``{{ name.var }}``. For complete information, you need to read the official one 
[documentation](https://github.com/flosch/pongo2).

__SetResponseWriter__
```
SetResponseWriter(w http.ResponseWriter)
```
Mandatory method. Sets the ``http.ResponseWriter`` to the templater.

__SetRequest__
```
SetRequest(r *http.Request)
```
Mandatory method. Sets ``*http.Request`` to the templater.

### Standard global variables are created in the framework
Global variables that are added using the framework are described here.

* csrf_token - a global variable that appears only when the parameter named ``csrf_token`` is set in cookies. Come on 
the variable must be used in the HTML form to validate the CSRF token.

## Filters
Brief description of filters created in the framework.

* unescape - descreens the text.
* strslice - converts a slice of any type into a string.

### Методи модуля фільтрів

__RegisterGlobalFilter__
```
RegisterGlobalFilter(name string, fn pongo2.FilterFunction) error
```
Registers a global filter.

__RegisterMultipleGlobalFilter__
```
RegisterMultipleGlobalFilter(filters []Filter) error
```
Registers multiple filters using the ``Filter`` structure.