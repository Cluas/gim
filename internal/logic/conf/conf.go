package conf

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/viper"
)

var (
	// Conf is config for logic server
	Conf     *Config
	confPath string
)

func init() {
	flag.StringVar(&confPath, "p", ".", " set logic config file path")
}

// Config is struct of logic config
type Config struct {
	Base  *BaseConf  `mapstructure:"base"`
	Redis *RedisConf `mapstructure:"redis"`
	RPC   *RPCConf   `mapstructure:"rpc"`
	HTTP  *HTTPConf  `mapstructure:"http"`
}

// BaseConf is struct of base config
type BaseConf struct {
	PidFile    string `mapstructure:"pidfile"`
	MaxProc    int
	PprofAddrs []string `mapstructure:"pprofbind"` // 性能监控的域名端口

}

// RedisConf is struct of redis config
type RedisConf struct {
	Password  string `mapstructure:"password"`
	DefaultDB int    `mapstructure:"default_db"`
	Address   string `mapstructure:"address"`
}

// RPCConf is config for logic rpc
type RPCConf struct {
	Address []string `mapstructure:"address"`
}

// HTTPConf is config for http server
type HTTPConf struct {
	Address      []string      `mapstructure:"address"`
	ReadTimeout  time.Duration `mapstructure:"HTTPReadTimeout"`
	WriteTimeout time.Duration `mapstructure:"HTTPWriteTimeout"`
}

// Init is func to intial logic config
func Init() (err error) {
	Conf = NewConfig()
	viper.SetConfigName("logic")
	viper.SetConfigType("toml")
	viper.AddConfigPath(confPath)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Errorf("Unable to decode into struct：%s", err))
	}

	return nil
}

// NewConfig is func to create a logic config
func NewConfig() *Config {
	return &Config{
		Base: &BaseConf{
			PidFile:    "/tmp/logic.pid",
			MaxProc:    runtime.NumCPU(),
			PprofAddrs: []string{"localhost:6971"},
		},
		Redis: &RedisConf{
			Password:  "redis123#",
			DefaultDB: 0,
			Address:   "localhost:6379",
		},
		RPC: &RPCConf{
			Address: []string{"tcp@localhost:6923"},
		},
		HTTP: &HTTPConf{
			Address:      []string{"tcp@0.0.0.0:6921"},
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 20 * time.Second,
		},
	}
}
