package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/uwine4850/foozy/pkg/utils/fmap"
	"gopkg.in/yaml.v3"
)

type Generate struct {
	config *Config
}

func NewGenerate(config *Config) *Generate {
	return &Generate{config: config}
}

func (g *Generate) SetConfig(config *Config) {
	g.config = config
}

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
