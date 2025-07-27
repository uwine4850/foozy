## package middlewares
This package implements middleware.<br>
Middleware is software that acts as a bridge between applications, operating systems, and other software components, providing common services and capabilities.

Example of use:
```golang
mdd := middlewares.NewMiddlewares()
mdd.PreMiddleware(0, func(w http.ResponseWriter, r *http.Request, m interfaces.Manager) error {
	fmt.Println("PRE")
	return nil
})
mdd.AsyncMiddleware(func(w http.ResponseWriter, r *http.Request, m interfaces.Manager) error {
	fmt.Println("ASYNC")
	return nil
})
mdd.AsyncMiddleware(func(w http.ResponseWriter, r *http.Request, m interfaces.Manager) error {
	fmt.Println("ASYNC 2")
	return nil
})
mdd.PostMiddleware(0, func(r *http.Request, m interfaces.Manager) error {
	fmt.Println("POST")
	return nil
})
adapter := router.NewAdapter(newManager, mdd)
```
__NOTE:__ you must install middlewares in the [adapter](/router/router/#adapter).
```golang
adapter := router.NewAdapter(newManager, mdd)
```

## Middlewares
Middlewares implementation of the middleware concept for the framework.
Three types of middleware are possible:

1. [PreMiddleware](#middlewarespremiddleware) — executed synchronously before the request. 
These middlewares are executed exactly according to their established order.
2. [PostMiddleware](#middlewarespostmiddleware) — executed synchronously after the request. These middlewares 
will be called in the specified order.
3. [AsyncMiddleware](#middlewaresasyncmiddleware) — asynchronous middleware. They are executed asynchronously, 
but after [PreMiddleware](#middlewarespremiddleware) and before the request handler. Cannot be ordered, so are not used in a chain.
Also each middleware returns an error, it is handled in the router. But if you need more specific processing, 
you can not return an error, but process it directly in the middleware.

Also very important point: there are also functions that can control the router from middleware, 
they are in the same package. These functions should be called only in middleware.
These functions include:

1. [SkipNextPage](#skipnextpage) — skips the page turner. That is, not just its display, but the whole logic.
2. [SkipNextPageAndRedirect](#skipnextpageandredirect) — does the same thing as [SkipNextPage](#skipnextpage), but does redirect after skipping the page.

#### Middlewares.PreMiddleware
Middleware that are executed in an ordered fashion before the url handler.<br>
The order must not be repeated.
```golang
func (mddl *Middlewares) PreMiddleware(order int, handler PreMiddleware) {
	if slices.Contains(mddl.preMiddlewaresOrder, order) {
		panic(fmt.Sprintf("middleware with order %s already exists", strconv.Itoa(order)))
	}
	mddl.preMiddlewaresOrder = append(mddl.preMiddlewaresOrder, order)
	mddl.preMiddlewares[order] = handler
}
```

#### Middlewares.PostMiddleware
Middleware that runs after the HTTP request handler has been confirmed.
```golang
func (mddl *Middlewares) PostMiddleware(order int, handler PostMiddleware) {
	if slices.Contains(mddl.postMiddlewaresOrder, order) {
		panic(fmt.Sprintf("middleware with order %s already exists", strconv.Itoa(order)))
	}
	mddl.postMiddlewaresOrder = append(mddl.postMiddlewaresOrder, order)
	mddl.postMiddlewares[order] = handler
}
```

#### Middlewares.AsyncMiddleware
Middleware that is executed asynchronously before the request handler, 
but after `PreMiddleware` processing.<br>
Can't create chains, not called in an orderly fashion.
```golang
func (mddl *Middlewares) AsyncMiddleware(handler AsyncMiddleware) {
	mddl.asyncMiddlewares = append(mddl.asyncMiddlewares, handler)
}
```

#### Middlewares.RunPreMiddlewares
Runs all `PreMiddleware`. Starts them in sorted order, i.e. 1...n.
```golang
func (mddl *Middlewares) RunPreMiddlewares(w http.ResponseWriter, r *http.Request, m interfaces.Manager) error {
	mddl.preMiddlewaresOrder.Sort()
	for i := 0; i < mddl.preMiddlewaresOrder.Len(); i++ {
		order := mddl.preMiddlewaresOrder[i]
		if err := mddl.preMiddlewares[order](w, r, m); err != nil {
			return err
		}
	}
	return nil
}
```

#### Middlewares. RunAndWaitAsyncMiddlewares
Runs asynchronous middlewares.<br>
It also waits for them to complete, no additional actions are needed.<br>
If at least one middleware causes an error, all handlers stop.
```golang
func (mddl *Middlewares) RunAndWaitAsyncMiddlewares(w http.ResponseWriter, r *http.Request, m interfaces.Manager) error {
	var wg sync.WaitGroup
	var asyncError error
	var mu sync.Mutex
	stop := make(chan struct{})
	for i := 0; i < len(mddl.asyncMiddlewares); i++ {
		handler := mddl.asyncMiddlewares[i]
		wg.Add(1)
		go func(h AsyncMiddleware) {
			defer wg.Done()

			// If at least one handler causes an error, all other handlers will fail to run.
			select {
			case <-stop:
				return
			default:
			}

			if err := h(w, r, m); err != nil {
				mu.Lock()
				asyncError = err
				close(stop)
				mu.Unlock()
			}
		}(handler)
	}
	wg.Wait()
	return asyncError
}
```

#### Middlewares.RunPostMiddlewares
Runs all `PostMiddleware`. Starts them in sorted order, i.e. 1...n.
```golang
func (mddl *Middlewares) RunPostMiddlewares(r *http.Request, m interfaces.Manager) error {
	mddl.postMiddlewaresOrder.Sort()
	for i := 0; i < mddl.postMiddlewaresOrder.Len(); i++ {
		order := mddl.postMiddlewaresOrder[i]
		if err := mddl.postMiddlewares[order](r, m); err != nil {
			return err
		}
	}
	return nil
}
```

#### SkipNextPage
Sends a command to the [router](/router/router) to skip rendering the next page.
```golang
func SkipNextPage(manager interfaces.ManagerOneTimeData) {
	manager.SetUserContext(namelib.ROUTER.SKIP_NEXT_PAGE, true)
	urlPattern, _ := manager.GetUserContext(namelib.ROUTER.URL_PATTERN)
	debug.RequestLogginIfEnable(debug.P_MIDDLEWARE, fmt.Sprintf("skip page at %s", urlPattern))
}
```

#### IsSkipNextPage
Checks if the page rendering should be skipped.
The function is built into the [router](/router/router).
```golang
func IsSkipNextPage(manager interfaces.ManagerOneTimeData) bool {
	_, ok := manager.GetUserContext(namelib.ROUTER.SKIP_NEXT_PAGE)
	return ok
}
```

#### SkipNextPageAndRedirect
Skips the page render and redirects to another page.
```golang
func SkipNextPageAndRedirect(manager interfaces.ManagerOneTimeData, w http.ResponseWriter, r *http.Request, path string) {
	http.Redirect(w, r, path, http.StatusFound)
	SkipNextPage(manager)
	debug.RequestLogginIfEnable(debug.P_MIDDLEWARE, fmt.Sprintf("redirect to %s", path))
}
```