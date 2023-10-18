## Package tmlengine
This package implements a templating engine based on the ``pongo2`` library.<br>
The package can be divided into two modules: the logic of the templating engine and the logic of filters for it.

## Template engine
The templating engine is used to easily display HTML code on a page in a browser. You can read more about pongo2
in the official repository at the following link: https://github.com/flosch/pongo2.

__SetPath__
```
SetPath(path string)
```
Required method. Sets the path to the HTML template.

__SetContext__
```
SetContext(data map[string]interface{})
```
Setting the context of the templating tool. Access to variables in HTML is carried out in the following way ``{{ name }}``. If the variables are <name>,
for example, a structure - access to variables will be as follows ``{{ name.var }}``. To get full information, you need to read the official
[documentation](https://github.com/flosch/pongo2).

__SetResponseWriter__
```
SetResponseWriter(w http.ResponseWriter)
```
Required method. Sets the ``http.ResponseWriter`` to the templating engine.

__SetRequest__
```
SetRequest(r *http.Request)
```
Required method. Sets the ``*http.Request`` to the template.

### Standard global variables created in the framework
This describes the global variables that are added using the framework.

* csrf_token - a global variable that appears only when the parameter named ``csrf_token`` is set to cookies. This
  variable should be used in the HTML form to validate the CSRF token.

## Фільтри
A brief description of the filters created in the framework.

* unescape - unescape the text.
* strslice - converts a slice of any type to a string.

### Methods of the filter module

__RegisterGlobalFilter__
```
RegisterGlobalFilter(name string, fn pongo2.FilterFunction) error
```
Registers a global filter.

__RegisterMultipleGlobalFilter__
```
RegisterMultipleGlobalFilter(filters []Filter) error
```
Registers several filters using the ``Filter'' structure.