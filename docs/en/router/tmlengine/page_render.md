## IRender
This interface is designed to simplify interaction with the template builder. He is also 
contains its own methods for displaying data on the page.

__New()__
```
New() (interface{}, error)
```
Implementation of the __INewInstance__ interface. Used for automatic 
creating a new instance. If the template generator is not set to manual, it will be 
the standard implementation of the template generator should be used.

__SetContext__
```
SetContext(data map[string]interface{})
```
Sets the context for the templater. In the template, it is possible to call the installed ones 
data by key.

__GetContext__
```
GetContext() map[string]interface{}
```
Returns the value of the templater context.

__SetTemplateEngine__
```
SetTemplateEngine(engine ITemplateEngine)
```
Installs the templater. No need to call to use standard implementation.

__GetTemplateEngine__
```
GetTemplateEngine() ITemplateEngine
```
Returns the templator in use.

__RenderTemplate__
```
RenderTemplate(w http.ResponseWriter, r *http.Request) error
```
Configures and runs the templater.

__SetTemplatePath__
```
SetTemplatePath(templatePath string)
```
Sets the path to the HTML template. Must be called before __RenderTemplate__.

__RenderJson__
```
RenderJson(data interface{}, w http.ResponseWriter) error
```
Displays data in JSON format on the page.

## Інші функції

__CreateAndSetNewRenderInstance__
```
CreateAndSetNewRenderInstance(manager interfaces.IManager) error
```
Creates and installs a new instance of the renderer in the manager. It is used in 
routers