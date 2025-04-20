package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/eclipse-xfsc/microservice-core-go/pkg/server"

	"github.com/spf13/viper"
)

var viperInstance = viper.New()

// BaseConfig can be used to import the base config parameters in the applications config struct.
// Please use with tag `mapstructure:",squash"`.
type BaseConfig struct {
	LogLevel   string            `mapstructure:"logLevel" envconfig:"LOG_LEVEL" default:"info"`
	IsDev      bool              `mapstructure:"isDev" envconfig:"IS_DEV" default:"false"`
	ListenAddr string            `mapstructure:"listenAddr" envconfig:"LISTEN_ADDR" default:"127.0.0.1"`
	ListenPort int               `mapstructure:"listenPort" envconfig:"LISTEN_PORT" default:"8080"`
	ServerMode server.ServerMode `mapstructure:"serverMode" default:"production"`
}

// LoadConfig sets given defaults and read in given config.
func LoadConfig(prefix string, config any, defaults map[string]any) error {
	setDefaults(defaults)

	if err := readConfig(prefix); err != nil {
		return err
	}

	if err := viperInstance.Unmarshal(config); err != nil {
		return err
	}

	return nil
}

func readConfig(prefix string) error {
	viperInstance.SetConfigName("config")
	viperInstance.SetConfigType("yaml")
	viperInstance.AddConfigPath(".")

	viperInstance.SetEnvPrefix(strings.ToTitle(prefix))
	viperInstance.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viperInstance.AutomaticEnv()

	if err := viperInstance.ReadInConfig(); err != nil {
		if !errors.Is(err, viper.ConfigFileNotFoundError{}) {
			return fmt.Errorf("error read in configFile: %w", err)
		}
	}

	return nil
}

func setDefaults(defaults map[string]any) {
	viperInstance.SetDefault("logLevel", "info")
	viperInstance.SetDefault("listenAddr", "127.0.0.1")
	viperInstance.SetDefault("listenPort", 8080)
	viperInstance.SetDefault("serverMode", server.ModeProduction)

	for key, value := range defaults {
		viperInstance.SetDefault(key, value)
	}
}
