package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/uwine4850/foozy/pkg/utils/fmap"
	"gopkg.in/yaml.v3"
)

// Generate structure for generating a .yaml configuration file.
// For proper operation, the Config object must be used.
// The .yaml file will be generated from it. In the standard implementation
// to get the Config object you should use the config.Cnf() method. It can be configured beforehand.
type Generate struct {
	config *Config
}

func NewGenerate(config *Config) *Generate {
	return &Generate{config: config}
}

// SetConfig sets the configuration file from which the .yaml file will be generated.
func (g *Generate) SetConfig(config *Config) {
	g.config = config
}

// Gen generates a configuration file. The previously installed Config object is used for generation.
// Before generation, the previous config file is loaded, if it exists. When the configuration file exists, the following actions are performed:
//
//	If the GeneratedDefault field is true, the default config will not be overwritten in the Config file.
//	If GeneratedAdditionally is true, the optional configuration is not overwritten, but if there are new or deleted fields, such changes will take effect.
func (g *Generate) Gen() error {
	loadConfig, err := Load()
	if err != nil && !reflect.DeepEqual(reflect.TypeOf(err), reflect.TypeOf(&ErrNoFile{})) {
		fmt.Println(reflect.TypeOf(err))
		return err
	}
	if loadConfig != nil {
		if loadConfig.GeneratedDefault {
			g.config.Default = loadConfig.Default
		}
		if loadConfig.GeneratedAdditionally {
			fmap.MergeMap(&g.config.Additionally, loadConfig.Additionally)
		}
	}
	file, err := os.Create(g.config.path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := yaml.NewEncoder(file)
	if err := encoder.Encode(g.config); err != nil {
		return err
	}
	return nil
}
