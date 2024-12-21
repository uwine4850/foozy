package config

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/uwine4850/foozy/pkg/typeopr"
)

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
