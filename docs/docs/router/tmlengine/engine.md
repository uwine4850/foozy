## TemplateEngine
An object that renders data on an HTML page. Implements the `TemplateEngine` interface.

The standard implementation uses the [pongo2](https://github.com/flosch/pongo2) template engine.

#### TemplateEngine.New
Implementation of the `NewInstance` interface.<br>
Creates a new instance of the `TemplateEngine` object with some predefined settings.
```golang
func (e *TemplateEngine) New() (interface{}, error) {
	RegisterMultipleGlobalFilter(BuiltinFilters)
	return &TemplateEngine{context: make(map[string]interface{})}, nil
}
```

#### TemplateEngine.SetPath
Sets the path to the HTML template.
```golang
func (e *TemplateEngine) SetPath(path string) {
	e.path = path
}
```

#### TemplateEngine.Exec
Does all the necessary processing for the template and shows the HTML code on the page.<br>
The CSRF token is also set in the context using the key `namelib.ROUTER.COOKIE_CSRF_TOKEN`. The token is taken from the cookie, and is usually set there using the __TODO:add link__ [SetCSRFToken]() function.
```golang
func (e *TemplateEngine) Exec() error {
	debug.RequestLogginIfEnable(debug.P_TEMPLATE_ENGINE, "exec template engine...")
	debug.RequestLogginIfEnable(debug.P_TEMPLATE_ENGINE, "processing html file")
	err := e.processingFile()
	if err != nil {
		return err
	}
	debug.RequestLogginIfEnable(debug.P_TEMPLATE_ENGINE, "set CSRF token")
	err = e.setCsrfVariable(e.request)
	if err != nil {
		return err
	}
	debug.RequestLogginIfEnable(debug.P_TEMPLATE_ENGINE, "execute template")
	execute, err := e.templateFile.Execute(e.context)
	if err != nil {
		return err
	}
	debug.RequestLogginIfEnable(debug.P_TEMPLATE_ENGINE, "write template")
	_, err = e.writer.Write([]byte(execute))
	if err != nil {
		return err
	}
	return nil
}
```

#### TemplateEngine.SetContext
Sets the context for the template engine. This context can be retrieved in the template using the following syntax: `{{ key }}`. If the template engine is passed to [pagerender](/router/tmlengine/pagerender) and it is in turn passed to [Manager](/router/manager/manager), then the context will always be empty for each new HTTP request.
```golang
func (e *TemplateEngine) SetContext(data map[string]interface{}) {
	fmap.MergeMapSync(&e.mu, &e.context, data)
}
```

#### TemplateEngine.GetContext
Getting context.
```golang
func (e *TemplateEngine) GetContext() map[string]interface{} {
	return e.context
}
```

#### TemplateEngine.SetResponseWriter
Set `http.ResponseWriter` for internal use.
```golang
func (e *TemplateEngine) SetResponseWriter(w http.ResponseWriter) {
	e.writer = w
}
```

#### TemplateEngine.SetRequest
Set `*http.Request` for internal use.
```golang
func (e *TemplateEngine) SetRequest(r *http.Request) {
	e.request = r
}
```

## Filters
Filters for the `pongo2` template engine.

#### Filter
An object that represents a single instance of a filter.
```golang
type Filter struct {
	Name string
	Fn   pongo2.FilterFunction
}
```

#### RegisterGlobalFilter
Registering global filters for templates.
```golang
func RegisterGlobalFilter(name string, fn pongo2.FilterFunction) {
	pongo2.RegisterFilter(name, fn)
}
```

#### RegisterMultipleGlobalFilter
Does the same thing as `RegisterGlobalFilter`, but accepts several filters at once as arguments.
```golang
func RegisterMultipleGlobalFilter(filters []Filter) {
	for i := 0; i < len(filters); i++ {
		RegisterGlobalFilter(filters[i].Name, filters[i].Fn)
	}
}
```