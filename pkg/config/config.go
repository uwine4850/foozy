package config

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/uwine4850/foozy/pkg/typeopr"
)

// Config is the main structure that describes the configuration of the framework.
// The GeneratedDefault and GeneratedAdditionally fields are intended to avoid re-generating the config
// and thus resetting the settings. If you still need to reset the settings, these fields should be set to false in the .yaml file.
// The Default field is responsible for the standard configuration of the framework.
// The Additionally field is for custom settings. That is, the user can put his own configs in this field
// and change them in the common .yaml file.
// The path and loadPath fields should contain the path to the generated configuration file.
// More precisely: path - the place of configuration generation, loadPath - the place of configuration
// loading. Two fields are made because the paths may differ syntactically, e.g. “config.yaml” and “../config.yaml”.
type Config struct {
	GeneratedDefault      bool                   `yaml:"GeneratedDefault"`
	GeneratedAdditionally bool                   `yaml:"GeneratedAdditionally"`
	Default               DefaultConfig          `yaml:"Config"`
	Additionally          map[string]interface{} `yaml:"Additionally"`
	path                  string
	loadPath              string
}

func (cnf *Config) SetPath(path string) {
	cnf.path = path
}

func (cnf *Config) SetLoadPath(path string) {
	cnf.loadPath = path
}

func (cnf *Config) AppendAdditionally(name string, item typeopr.IPtr) {
	cnf.Additionally[name] = item.Ptr()
}

var instance *Config
var once sync.Once

// Cnf singleton to access configuration settings.
// This function has nothing to do with outputting the configuration from
// the .yaml file, this function only provides access to the config generation settings.
func Cnf() *Config {
	once.Do(func() {
		instance = &Config{
			GeneratedDefault:      true,
			GeneratedAdditionally: true,
			Additionally:          map[string]interface{}{},
			Default: DefaultConfig{
				Debug: DebugConfig{
					PrintInfo:          true,
					Debug:              true,
					SkipLoggingLevel:   3,
					ErrorLoggingPath:   "errors.log",
					RequestInfoLogPath: "request.log",
				},
			},
		}
	})
	return instance
}

// DefaultConfig a standard set of configurations.
type DefaultConfig struct {
	Debug DebugConfig `yaml:"Debug"`
}

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

// Info displays information about each command.
// To work, you need to use the "i" tag. If this tag is missing, the command will be ignored.
func Info() {
	cnf := Cnf()
	defaultConfigType := reflect.TypeOf(cnf.Default)
	cnfObjectInfo("", &defaultConfigType)
	if len(cnf.Additionally) != 0 {
		fmt.Println("Additionally:")
		for _, configObject := range cnf.Additionally {
			configObjectType := reflect.TypeOf(configObject).Elem()
			cnfObjectInfo(" ", &configObjectType)
		}
	}
}

// cnfObjectInfo information about each command in the structure.
func cnfObjectInfo(sep string, object *reflect.Type) {
	objectType := *object
	fmt.Println(sep + objectType.Name() + ":")
	for i := 0; i < objectType.NumField(); i++ {
		if objectType.Field(i).Type.Kind() == reflect.Struct {
			ty := objectType.Field(i).Type
			cnfObjectInfo(sep+" ", &ty)
		} else {
			tag := objectType.Field(i).Tag
			yamlName := tag.Get("yaml")
			info := tag.Get("i")
			if yamlName != "" && info != "" {
				fmt.Println(sep + " " + yamlName + " — " + info)
			}
		}
	}
}
