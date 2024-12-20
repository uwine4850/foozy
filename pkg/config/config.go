package config

import (
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
	PrintInfo             bool   `yaml:"PrintInfo"`
	Debug                 bool   `yaml:"Debug"`
	DebugRelativeFilepath bool   `yaml:"DebugRelativeFilepath"`
	ErrorLogging          bool   `yaml:"ErrorLogging"`
	ErrorLoggingPath      string `yaml:"ErrorLoggingPath"`
	RequestInfoLog        bool   `yaml:"RequestInfoLog"`
	RequestInfoLogPath    string `yaml:"RequestInfoLogPath"`
	SkipLoggingLevel      int    `yaml:"SkipLoggingLevel"`
}
