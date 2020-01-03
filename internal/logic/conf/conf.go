package conf

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/viper"

	"github.com/Cluas/gim/pkg/log"
)

var (
	// Conf is config for logic server
	Conf     *Config
	confFile string
)

func init() {
	flag.StringVar(&confFile, "c", ".", " set logic config file.")
}

// Config is struct of logic config
type Config struct {
	Base  *BaseConf   `mapstructure:"base"`
	Redis *RedisConf  `mapstructure:"redis"`
	RPC   *RPCConf    `mapstructure:"rpc"`
	HTTP  *HTTPConf   `mapstructure:"http"`
	Log   *log.Config `mapstructure:"log"`
}

// BaseConf is struct of base config
type BaseConf struct {
	PidFile   string   `mapstructure:"pid_file"`
	PprofBind []string `mapstructure:"pprof_bind"`
	MaxProc   int
}

// RedisConf is struct of redis config
type RedisConf struct {
	Password    string        `mapstructure:"password"`
	Address     string        `mapstructure:"address"`
	DefaultDB   int           `mapstructure:"default_db"`
	MaxRetries  int           `mapstructure:"max_retries"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
}

// RPCConf is config for logic rpc
type RPCConf struct {
	Address []string `mapstructure:"address"`
}

// HTTPConf is config for http server
type HTTPConf struct {
	Address           []string      `mapstructure:"address"`
	ReadTimeout       time.Duration `mapstructure:"read_timeout"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
	IdleTimeout       time.Duration `mapstructure:"idle_timeout"`
}

// Init is func to initial logic config
func Init() (err error) {
	Conf = NewConfig()
	viper.SetConfigFile(confFile)

	if err = viper.ReadInConfig(); err != nil {
		return err
	}

	if err = viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Errorf("Unable to decode into struct：%s ", err))
	}

	return nil
}

// NewConfig is func to create a logic config
func NewConfig() *Config {
	return &Config{
		Base: &BaseConf{
			PidFile:   "/tmp/logic.pid",
			MaxProc:   runtime.NumCPU(),
			PprofBind: []string{"localhost:6971"},
		},
		Log: &log.Config{
			ServiceName: "logic",
		},
		Redis: &RedisConf{
			Password:    "redis123#",
			Address:     "localhost:6379",
			DefaultDB:   0,
			MaxRetries:  3,
			IdleTimeout: 5 * time.Second,
		},
		RPC: &RPCConf{
			Address: []string{"tcp@localhost:6923"},
		},
		HTTP: &HTTPConf{
			Address:           []string{"tcp@0.0.0.0:6921"},
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      20 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,
			IdleTimeout:       10 * time.Second,
		},
	}
}
