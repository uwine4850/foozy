## cmd

This package is used to run console commands.

To start working with this package, you need to create a file with any name you want, but for it to be the main package. Then create a function `main`, it can look like this:
```
func main() {
    if err := cmd.Run(); err != nil {
		panic(err)
    }
}
```
After that, you need to initialize the configuration settings. To do this, use the same function with the following command `go run <path to cmd.go> initcnf <target dir>`. This command creates a file that contains a configuration setting function, [read more here](https://github.com/uwine4850/foozy/blob/master/docs/config/config.md). It should be used in the `main` function. Now the function can look like the following:
```
func main() {
	initcnf.InitCnf()
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
```
After these operations you can fully use cmd.

The package contains the following commands:
* *cnfinfo* — configuration field information. Only information from the `i` tag is output.
* *initcnf <target_dir>* — generates a file with configuration settings.
* *gencnf* — generates a configuration. It is important that the configuration is pre-installed, e.g. with `initcnf`.

__Run__
```
Run() error
```
Starts the execution of commands.