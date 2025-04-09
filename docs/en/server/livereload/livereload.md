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
	initcnf "github.com/uwine4850/foozy/mycmd/init_cnf"
	"github.com/uwine4850/foozy/pkg/server/livereload"
)

func main() {
	initcnf.InitCnf()
	wrt := livereload.NewWiretap()
	wrt.SetDirs([]string{"dir"})
	reload := livereload.NewReloader("main.go", wrt)
	if err := reload.Start(); err != nil{
		panic(err)
	}
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
Sets the directories to be listened to.
It is important to specify that it is the directory and all files in it that is listened to.
One directory and all files in it is one `ObservedElement`.
Subdirectories are already considered new `ObservedElement`.

__SetExcludeDirs__
```
SetExcludeDirs(dirs []string)
```
Excludes the directory and absolutely all subdirectories from listening.

__SetFiles__
```
SetFiles(files []string)
```
Adds individual files to the wiretap.

__OnStart__
```
OnStart(fn func())
```
Starts every time during the start of the wiretapping.

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
element that was saved.

__Start__
```
Start() error
```
Starts listening.

### Example of use
```go
wiretap := livereload.NewWiretap()
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