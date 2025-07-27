## Config
This package is designed to create, customize and use a configuration file. The configuration uses the _singleton_ pattern. Also the configuration must be preloaded using the generated `InitCnf()` function. This function does not need to be called more than once per session.

__NOTE:__ after changing the configuration, you must reboot the server.

This package is divided into three sections:

* [Representing the configuration as a golang object](#representing-the-configuration-as-a-golang-object)
* [Generating a configuration file](#generating-a-configuration-file)
* [Loading the configuration file](#loading-the-configuration-file)

You can read more about each section below.

### Representing the configuration as a golang object
The `Config` object is used for the global representation of the configuration. The object stores configuration data, is used for generation and stores loaded data from `config.yaml`.
```golang
type Config struct {
	GeneratedDefault      bool                   `yaml:"GeneratedDefault"`
	GeneratedAdditionally bool                   `yaml:"GeneratedAdditionally"`
	Default               DefaultConfig          `yaml:"Config"`
	Additionally          map[string]interface{} `yaml:"Additionally"`
	path                  string
	loadPath              string
}
```
* The `GeneratedDefault` and `GeneratedAdditionally` fields are intended to avoid re-generating the config 
and thus resetting the settings. If you still need to reset the settings, these fields should be set to false in the config.yaml file.
* The `Default` field is responsible for the standard configuration of the framework.
* The `Additionally` field is for custom settings. That is, the user can put his own configs in this field and change them in the 
common config.yaml file. This should be done right here using the `AppendAdditionally` method:
```golang
func InitCnf() {
	cnf := config.Cnf()
	// Use this to add your configurations.
	// cnf.AppendAdditionally("my_cnf", typeopr.Ptr{}.New(&MyCnfCtruct{}))
	cnf.SetPath("cnf/config.yaml")
	cnf.SetLoadPath("cnf/config.yaml")
}
```
* The `path` and `loadPath` fields should contain the path to the generated configuration file. 

	More precisely: `path` - the place of configuration generation, `loadPath` - the place of configuration
	loading. Two fields are made because the paths may differ syntactically, e.g. "config.yaml" and "../config.yaml".

---
The `Cnf` function is a singleton for accessing configuration settings.
This function has nothing to do with outputting the configuration from the config.yaml file, it only provides access to the configuration generation settings. To just get the default configuration template you need to call this function.

__NOTE:__ this function should not be confused with `LoadedConfig` as it gives the same object but it does not store the actual project settings.
```golang
func Cnf() *Config
```
---
#### Configuration templates

All configuration templates are shown below. They will be used to generate the standard `config.yaml`. This is also an example of how the configuration passed to the `AppendAdditionally` method should look like.

Two structural tags will also be used here:

* yaml - field name.
* i - brief information about the field.

The `DefaultConfig` object is a master template that contains all additional configuration templates.
```golang
type DefaultConfig struct {
	Debug    DebugConfig    `yaml:"Debug"`
	Database DatabaseConfig `yaml:"Database"`
}
```

The `DebugConfig` object is the debug settings.
```golang
type DebugConfig struct {
	PrintInfo             bool   `yaml:"PrintInfo" i:"Displays basic information about each request."`
	Debug                 bool   `yaml:"Debug" i:"Enables debugging"`
	DebugRelativeFilepath bool   `yaml:"DebugRelativeFilepath" i:"In logs, file paths are displayed relatively"`
	ErrorLogging          bool   `yaml:"ErrorLogging" i:"Enables error logging"`
	ErrorLoggingPath      string `yaml:"ErrorLoggingPath" i:"Path to error log file"`
	RequestInfoLog        bool   `yaml:"RequestInfoLog" i:"Enables request logging"`
	RequestInfoLogPath    string `yaml:"RequestInfoLogPath" i:"Path to request log file"`
	SkipLoggingLevel      int    `yaml:"SkipLoggingLevel" i:"Skips logging levels. May need to be configured per project"`
}
```

The `DatabaseConfig` object is a part of the database configuration. The rest of the settings can be found directly in the __TODO: add link__ [corresponding package]().
```golang
type DatabaseConfig struct {
	MainConnectionPoolName string `yaml:"MainConnectionPoolName" i:"The name of the main connection pool"`
}
```
---
#### Info
The `Info` function displays all information about the configuration fields that have the `i` tag.
```golang
func Info()
```

### Generating a configuration file
There is only one `Generate` object in this section. It is used only for generating the `config.yaml` configuration file.
For correct operation it is necessary to use the `Config` object. A `config.yaml` file will be generated based on it. In the standard implementation, the `config.Cnf()` method must be used to get the `Config` object. It can be configured in advance.
---
#### Config.Gen
Method `Gen` generates a configuration file. The previously installed `Config` object is used for generation.
Before generation, the previous config file is loaded, if it exists. When the configuration file exists, the following actions are performed:

* If the GeneratedDefault field is true, the default config will not be overwritten in the Config file
* If GeneratedAdditionally is true, the optional configuration is not overwritten, but if there are new or deleted fields, such changes will take effect.
```golang
func (g *Generate) Gen() error
```
---
#### InitCnf
<a id="init-cnf"></a>
The `InitCnf` function is designed to initialize the nonfiguration. It currently does two things:

* Setting paths to configuration files.
* Adding a custom configuration via the `AppendAdditionally` method.

Therefore, this function must be called at least once per session in order to know where to load the configuration and also to load the user data.

```golang
func InitCnf() {
	cnf := config.Cnf()
	// Use this to add your configurations.
	// cnf.AppendAdditionally("my_cnf", typeopr.Ptr{}.New(&MyCnfCtruct{}))
	cnf.SetPath("cnf/config.yaml")
	cnf.SetLoadPath("cnf/config.yaml")
}
```
### Loading the configuration file
This part of the package is designed to load a sonfiguration from the `config.yaml` file and retrieve the loaded configuration.

#### Load
Function `Load` loads the settings from the `config.yaml` file. The file is loaded by the `path(loadPath)`, which is configured with the `Config` object. In addition, the configurations are loaded into a new instance of the `Config` structure.
Due to the fact that the field `Additionally` has type map[string]interface{} type interface{}
needs to be converted into a structure.
```golang
func Load() (*Config, error)
```
---
#### LoadedConfig
Function `LoadedConfig` singleton that loads the settings. It is preferable to use this function rather than `Load` to avoid loading the config every time.
It is important to specify that changes in `config.yaml` settings are not taken into account after the framework is started, so it is necessary to restart the server after changing the configuration.
```golang
func LoadedConfig() *Config
```