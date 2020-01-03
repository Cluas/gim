package conf

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/spf13/viper"

	"github.com/Cluas/gim/pkg/log"
)

var (
	// Conf is  var of job config
	Conf     *Config
	confFile string
)

func init() {
	flag.StringVar(&confFile, "c", ".", " set job config file.")
}

// Config is struct of config struct
type Config struct {
	Base  *BaseConf   `mapstructure:"base"`
	Log   *log.Config `mapstructure:"log"`
	Redis *RedisConf  `mapstructure:"redis"`
	Comet []CometConf `mapstructure:"comet"`
}

// RedisConf is struct of redis config
type RedisConf struct {
	Address   string `mapstructure:"address"`
	Password  string `mapstructure:"password"`
	DefaultDB int    `mapstructure:"default_db"`
}

// BaseConf is struct of base config
type BaseConf struct {
	PidFile      string   `mapstructure:"pid_file"`
	PprofBind    []string `mapstructure:"pprof_bind"`
	PushChan     int      `mapstructure:"push_chan"`
	PushChanSize int      `mapstructure:"push_chan_size"`
	IsDebug      bool
	MaxProc      int
}

// CometConf is struct of comet RPC
type CometConf struct {
	Key  int8   `mapstructure:"key"`
	Addr string `mapstructure:"addr"`
}

// Init is func to initial log config
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

// NewConfig is func to create a Config
func NewConfig() *Config {
	return &Config{
		Base: &BaseConf{
			PidFile:      "/tmp/job.pid",
			MaxProc:      runtime.NumCPU(),
			PprofBind:    []string{"localhost:6922"},
			PushChan:     2,
			PushChanSize: 50,
			IsDebug:      true,
		},
		Redis: &RedisConf{
			Address:   "127.0.0.1:6379",
			Password:  "redis123#",
			DefaultDB: 0,
		},
		Comet: []CometConf{
			{Key: 1, Addr: "tcp@0.0.0.0:6999"},
		},
		Log: &log.Config{
			ServiceName: "job",
		},
	}
}
