package config

import (
	"flag"
	"fmt"

	"github.com/Cluas/gim/pkg/log"
	"github.com/spf13/viper"
)

// Config is struct of comet config
type Config struct {
	Base *BaseConfig `mapstructure:"base"`
	Log  *log.Config `mapstructure:"log"`
}

// BaseConfig is struct of base config
type BaseConfig struct {
	PidFile string `mapstructure:"pidfile"`
}

var (
	// Conf is var of config
	Conf       *Config
	configPath string
)

func init() {
	flag.StringVar(&configPath, "d", "./internal/comet/config/", "set logic config file path")
}

// Init is func to initial config
func Init() (err error) {
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
		Base: &BaseConfig{
			PidFile: "/tmp/comet.pid",
		},
		Log: &log.Config{
			LogPath:  "/Users/cluas/code/gim/log.log",
			LogLevel: "debug",
		},
	}
}
