## Config
This package is designed to interact with the configuration. With this package you 
can both generate and read the configuration.

The Config object is used for two things:
* Setting up configuration generation. That is, the configuration will be generated 
according to the structure of this object. Therefore, you can fill the required fields 
of the object with data before generation, such fields will have default 
values. You can also use the `Additionally` field to put custom settings there.
* Loading data from the configuration. The data that the `Load()` function loads 
places the data into an instance of the `Config` object.

To get started with the configuration, you must first set up 
generation. You can use cmd command `go run mycmd/cmd.go cnf-init <target_dir>` to configure it. This 
command uses codegeneration to create basic settings. [This file](https://github.com/uwine4850/foozy/blob/master/internal/codegen/init_cnf/init_cnf.go) 
will be generated.

The `InitCnf()` function will be used as a configuration setting. This function can be freely 
changed for the desired configuration setting. The function is used in several places:
* Before initializing the [Cmd](https://github.com/uwine4850/foozy/blob/master/docs/cmd/cmd.md) object. 
This is necessary for the configuration generation to be successful.
* At the beginning of the `main()` function. This is necessary for the configuration loading to work properly.
* Before using [livereloader](https://github.com/uwine4850/foozy/blob/master/docs/en/server/livereload/livereload.md).

When everything is set up, you can proceed to configuration generation using the initialized 
[Cmd](https://github.com/uwine4850/foozy/blob/master/docs/cmd/cmd.md) object. To do this, 
use the `go run mycmd/cmd.go cnf-gen` command. After the command is executed, the config will be generated.

In the default implementation, it is not possible to retrieve a configuration value 
if it has changed since the server was started. But if you put the directory with 
the configuration file in livereload, you can reboot the server after each 
configuration save, and thus retrieve the changes.

It is important to mention the Debug configuration field. It is used in the following places:
* Output more detailed errors.
* Disabling and enabling the reloader.

__Cnf__
```
Cnf() *Config
```
Singleton for accessing configuration settings. It is with this function that the configuration must be set up.

__Info__
```
Info()
```
Outputs information about all configuration fields only if there is an “i” tag there.

__Generate.Gen__
```
Gen() error
```
Generates a .yaml configuration file.

__Load__
```
Load() (*Config, error)
```
Loads the configuration from a .yaml file.

__LoadedConfig__
```
LoadedConfig() *Config
```
Provides access to the singleton of the loaded config. It is desirable to use this 
singleton to access the configuration, since it loads it once. Also, bypassing it, 
the configuration will not be updated in real time, but only after a server restart.