package config

import (
	"fmt"
	"os"
	"reflect"
	"sync"

	"github.com/uwine4850/foozy/pkg/typeopr"
	"github.com/uwine4850/foozy/pkg/utils/fmap"
	"github.com/uwine4850/foozy/pkg/utils/fstring"
	"gopkg.in/yaml.v3"
)

// Load loads the settings from the .yaml file. The file is loaded by the path(loadPath),
// which is configured with the Config object. In addition, the configurations are loaded
// into a new instance of the Config structure.
// Due to the fact that the field In addition has type map[string]interface{} type interface{}
// needs to be converted into a structure. For this purpose, we need to use the Config singleton
// accessed by the Cnf() function.
func Load() (*Config, error) {
	loadPath := Cnf().loadPath
	if !fstring.PathExist(loadPath) {
		return nil, &ErrNoFile{Path: loadPath}
	}

	var config Config

	file, err := os.Open(loadPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	newAdditionallyMap := make(map[string]interface{}, len(config.Additionally))
	for name, value := range config.Additionally {
		additionallyObject, ok := Cnf().Additionally[name]
		if ok {
			var additionallyObjectType reflect.Type
			if reflect.TypeOf(additionallyObject).Kind() == reflect.Ptr {
				additionallyObjectType = reflect.TypeOf(additionallyObject).Elem()
			} else {
				additionallyObjectType = reflect.TypeOf(additionallyObject)
			}
			if additionallyObjectType.Kind() == reflect.Struct {
				newAdditionallyStruct := reflect.New(additionallyObjectType).Elem()
				vv := value.(map[string]interface{})
				if err := fmap.YamlMapToStruct(&vv, typeopr.Ptr{}.New(&newAdditionallyStruct)); err != nil {
					return nil, err
				}
				newAdditionallyMap[name] = newAdditionallyStruct.Addr().Interface()
			} else {
				newAdditionallyMap[name] = value
			}
		}
	}
	config.Additionally = newAdditionallyMap
	return &config, nil
}

type ErrNoFile struct {
	Path string
}

func (e *ErrNoFile) Error() string {
	return fmt.Sprintf("config file %s not exists", e.Path)
}

var loadedConfigInstance *Config
var onceLoadedConfig sync.Once

// LoadedConfig singleton that loads the settings. It is preferable
// to use this function rather than Load to avoid loading the config every time.
// It is important to specify that changes in .yaml settings are not taken into
// account after the framework is started, so it is necessary to restart the server
// after changing the configuration.
func LoadedConfig() *Config {
	onceLoadedConfig.Do(func() {
		config, err := Load()
		if err != nil {
			panic(err)
		}
		loadedConfigInstance = config
	})
	return loadedConfigInstance
}
