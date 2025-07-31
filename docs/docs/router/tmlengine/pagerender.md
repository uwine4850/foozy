## pagerender
This module is designed for rendering HTML pages using a template engine. An object that implements the `TemplateEngine` interface is used for rendering.

## Render
The object that controls the installed template engine.
This object implements the `Render` and `NewInstance` interfaces.

Example of use:
```golang
...
newRender, err := tmlengine.NewRender()
if err != nil {
		panic(err)
}
newManager := manager.NewManager(
	manager.NewOneTimeData(),
	newRender,
	database.NewDatabasePool(),
)
...
newRouter.Register(router.MethodGET, "/page",
	func(w http.ResponseWriter, r *http.Request, manager interfaces.Manager) error {
		manager.Render().SetTemplatePath("index.html")
		if err := manager.Render().RenderTemplate(w, r); err != nil {
			return err
		}
		return nil
	})
```
#### Render.SetContext
Setting the context of the template engine.
```golang
func (rn *Render) SetContext(data map[string]interface{}) {
	rn.TemplateEngine.SetContext(data)
}
```

#### Render.GetContext
Getting context from the template engine.
```golang
func (rn *Render) GetContext() map[string]interface{} {
	return rn.TemplateEngine.GetContext()
}
```
#### Render.SetTemplateEngine
Set the template engine interface.
Optional method if the template engine is already installed.
```golang
func (rn *Render) SetTemplateEngine(engine interfaces.TemplateEngine) {
	rn.TemplateEngine = engine
}
```

#### Render.RenderTemplate
Calls the methods of the `TemplateEngine` interface to configure and render the HTML template.
```golang
func (rn *Render) RenderTemplate(w http.ResponseWriter, r *http.Request) error {
	if rn.templatePath == "" {
		return ErrTemplatePathNotSet{}
	}
	if !fpath.PathExist(rn.templatePath) {
		return ErrTemplatePathNotExist{Path: rn.templatePath}
	}
	rn.TemplateEngine.SetPath(rn.templatePath)
	rn.TemplateEngine.SetResponseWriter(w)
	rn.TemplateEngine.SetRequest(r)
	err := rn.TemplateEngine.Exec()
	if err != nil {
		return err
	}
	return nil
}
```

#### Render.SetTemplatePath
Setting the path to the template that the templating engine renders.
```golang
func (rn *Render) SetTemplatePath(templatePath string) {
	rn.templatePath = templatePath
}
```