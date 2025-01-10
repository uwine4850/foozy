## Package manager
This package implements the manager's work algorithm. A manager is needed for some
project settings and to transfer data from the server to the http request handler.
At the time of writing, the manager implements four interfaces that perform the following 
occupation:<br>
* interfaces.IManagerOneTimeData — stores one-time data for each request.
* interfaces.IRender — used to display the html page.
* interfaces.IKey — сontrols the keys.

It is important to note that IManagerOneTimeData and IRender are different from others,
 because they are unique for each request.

Implementations of all these interfaces are found in the ``Manager`` structure, which looks like this
 as follows:
```
type Manager struct {
    managerData      interfaces.IManagerOneTimeData
    render      interfaces.IRender
	key         interfaces.IKey
}
```
In this structure there are simple "get" and "set" methods to implement and 
getting individual managers.

It is important to mention that `IRender` does not have to be installed. 
If it is not installed, it simply will not be accessible. All the functionality 
that sets the context in the renderer duplicates this data in `UserContext`. 
This does not apply to `object.TemplateView`, this object should always use IRender as it is meaningless without it.

You can read more about each of the managers at the link:

* [IManagerOneTimeData](https://github.com/uwine4850/foozy/blob/master/docs/en/router/manager/manager_otd.md)
* [IRender](https://github.com/uwine4850/foozy/blob/master/docs/en/router/tmlengine/page_render.md)
* [IKey](https://github.com/uwine4850/foozy/blob/master/docs/en/secure/key.md)