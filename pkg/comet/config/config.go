package config

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

// Config is struct of comet config
type Config struct {
	Base BaseConf `mapstructure:"base"`
}

// BaseConf is struct of base config
type BaseConf struct {
	pidfile string `mapstructure:"pidfile"`
}

var (
	// Conf is var of config
	Conf       *Config
	configPath string
)

func init() {
	flag.StringVar(&configPath, "d", "./", "set logic config file path")
}

// InitConfig is func to initial config
func InitConfig() (err error) {
	Conf = NewConfig()
	viper.SetConfigName("comet")
	viper.SetConfigType("toml")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Errorf("unable to decode into struct:  %s \n ", err))
	}
	return nil
}

// NewConfig is constructor of Conig
func NewConfig() *Config {
	return &Config{
		Base: BaseConf{
			pidfile: "/tmp/comet.pid",
		},
	}
}
