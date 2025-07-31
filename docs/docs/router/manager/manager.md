## Manager
The `Manager` object is used to manage router modules. The object's responsibilities include the following tasks:

* Transferring data within a single HTTP request. That is, data cannot be transferred between two HTTP requests
* Storing the [Render](/router/tmlengine/tmlengine) object
* More convenient access to the __TODO: link__ [Key]() module
* Access to __TODO: link__ [DatabasePool]()

#### Manager.New
Creates a new instance of `Manager` with some of the old settings. This is a very important method because it creates 
a new instance of `Manager`, but retains the static data. If a new instance is created without this method, the [router](/router/router) 
may not work properly or may cause serious problems. Also creates a new instance of [Render](/router/tmlengine/tmlengine) if it has been set previously.
```golang
func (m *Manager) New() (interface{}, error) {
	newOTD, err := m.oneTimeData.New()
	if err != nil {
		return nil, err
	}
	var newRender interfaces.Render
	if m.render != nil {
		_newRender, err := m.render.New()
		if err != nil {
			return nil, err
		}
		newRender = _newRender.(interfaces.Render)
	}

	return &Manager{
		oneTimeData:  newOTD.(interfaces.ManagerOneTimeData),
		render:       newRender,
		key:          m.key,
		databasePool: m.databasePool,
	}, nil
}
```

## OneTimeData
An object that transfers data between router modules. Data can only be transferred within the boundaries of a single HTTP request. For example, data cannot be transferred to another handler. For correct operation, a new instance of this object must be created in [Adapter](/router/router/#adapter) for each handler call.

#### OneTimeData.SetSlugParams
Sets slug parameters for further retrieval by the user. The setting is made inside [adapter](/router/router/#adapter).
```golang
func (m *OneTimeData) SetSlugParams(params map[string]string) {
	m.slugParams = params
}
```

#### OneTimeData.GetSlugParams
Returns the slug parameter.
```golang
func (m *OneTimeData) GetSlugParams(key string) (string, bool) {
	res, ok := m.slugParams[key]
	return res, ok
}
```

#### OneTimeData.SetUserContext
Sets the user context. The user can then use this data. The framework automatically sets some data here, here is a list of it:

* `namelib.ROUTER.URL_PATTERN` — the current URL pattern.
* `namelib.ROUTER.REDIRECT_ERROR` — redirect error. Set only if the [router.CatchRedirectError](/router/router/#catchredirecterror) function is called.
* `namelib.ROUTER.SERVER_ERROR` — server error. Set only if the [router.ServerError](/router/router/#servererror) function is called.
* `namelib.ROUTER.SERVER_FORBIDDEN_ERROR` — access error. Set only if the [router.ServerForbidden](/router/router/#serverforbidden) function is called.
* `namelib.ROUTER.SKIP_NEXT_PAGE` — tells the router to skip the page handler. Set only if the __TODO: link__ [middlewares.SkipNextPage]() function is called.
* `namelib.OBJECT.OBJECT_CONTEXT` — object that is filled in __TODO: link__ [view]().
* `namelib.ROUTER.COOKIE_CSRF_TOKEN` — html string with CSRF token. Set only if the __TODO: link__ [secure.SetCSRFToken]() function is called.
```golang
func (m *OneTimeData) SetUserContext(key string, value interface{}) {
	m.userContext.Store(key, value)
}
```

#### OneTimeData.GetUserContext
Returns the user context.
```golang
func (m *OneTimeData) GetUserContext(key string) (any, bool) {
	m.userContext.Range(func(key, value any) bool {
		return true
	})
	value, ok := m.userContext.Load(key)
	return value, ok
}
```

#### OneTimeData.DelUserContext
Removes user context.
```golang
func (m *OneTimeData) DelUserContext(key string) {
	m.userContext.Delete(key)
}
```