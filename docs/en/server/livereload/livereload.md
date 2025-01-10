## Package livereload
This package is designed to simplify and speed up development and should only be used during the development phase.<br>
When changes are made to project files, the project needs to be constantly reloaded. This package is needed to solve 
this problem and do it automatically.<br>
This package is divided into two parts, namely the part with rebooting and the part with listening
saving files.

### Example of use
```go
package main

import (
    "github.com/uwine4850/foozy/pkg/livereload"
    initcnf "github.com/uwine4850/foozy/mycmd/init_cnf"
)

func main() {
    initcnf.InitCnf()
	reload := livereload.NewReload("project/cmd/main.go", livereload.NewWiretap([]string{"project", "pkg"},
		[]string{}))
	reload.Start()
}
```
In this example, the server that is located in the file ``project/cmd/main.go`` is overloaded. A reboot occurs
when any file is saved in the ``project`` or ``pkg`` directory.


## Listen to save files
The ``IWiretap`` interface is responsible for this functionality. Next, we will describe the methods that are associated with it.
__SetDirs__
```
SetDirs(dirs []string)
```
The required method. It is used to set the directories in which files will be listened to.

__OnStart__
```
OnStart(fn func())
```
The method runs the function once ``fn`` at the start of listening.

__GetOnStartFunc__
```
GetOnStartFunc() func()
```
Returns the function that was set by the __OnStart__ method.

__OnTrigger__
```
OnTrigger(fn func(filePath string))
```
The method sets the ``fn`` function that is executed every time the file is saved. The ``filePath string`` parameter is the path to the
file that was saved.

__SetUserParams__
```
SetUserParams(key string, value interface{})
```
Sets the parameters that can be passed between the functions ``OnStart(fn func())`` and ``OnTrigger(fn func(filePath string))``.

__GetUserParams__
```
GetUserParams(key string) (interface{}, bool)
```
Returns the user parameters set by the __SetUserParams__ method.

__Start__
```
Start() error
```
Starts listening.

### Example of use
```go
wiretap := livereload.NewWiretap3()
wiretap.SetDirs([]string{"project", "project_files"})
wiretap.OnStart(func() {
    fmt.Println("Start.")
})
wiretap.OnTrigger(func(filePath string) {
    fmt.Println("Trigger.")
})
err := wiretap.Start()
if err != nil {
    panic(err)
}
```

## Restarting the server
To implement this functionality, the ``Reload`` structure is used.

Constructor __NewReload(pathToServerFile string, wiretap interfaces.IWiretap) *Reload__<br>
* pathToServerFile - the path to the file that starts the server, for example, ``project/cmd/main.go``.
* wiretap - an instance of ``interfaces.IWiretap``.

__Start__
```go
Start()
```
Starting a server reboot.