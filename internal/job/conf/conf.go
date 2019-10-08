package conf

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/spf13/viper"
)

var (
	Conf     *Config
	confPath string
)

func init() {
	flag.StringVar(&confPath, "p", ".", " set job config file path")
}

type Config struct {
	Base  *BaseConf   `mapstructure:"base"`
	Redis *RedisConf  `mapstructure:"redis"`
	Comet []CometConf `mapstructure:"comet"`
	// Bucket BucketConf `mapstructure:"bucket"`
}

type RedisConf struct {
	Address   string `mapstructure:"address"` //
	Password  string `mapstructure:"password"`
	DefaultDB int    `mapstructure:"default_db"`
}

// 基础的配置信息
type BaseConf struct {
	Pidfile      string `mapstructure:"pidfile"`
	MaxProc      int
	PprofAddrs   []string `mapstructure:"pprofBind"`
	PushChan     int      `mapstructure:"pushChan"`
	PushChanSize int      `mapstructure:"pushChanSize"`
	IsDebug      bool
}
type CometConf struct {
	Key  int8   `mapstructure:"key"`
	Addr string `mapstructure:"addr"`
}

func Init() (err error) {
	Conf = NewConfig()
	viper.SetConfigName("job")
	viper.SetConfigType("toml")
	viper.AddConfigPath(confPath)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Errorf("unable to decode into struct：  %s \n", err))
	}

	return nil
}

func NewConfig() *Config {
	return &Config{
		Base: &BaseConf{
			Pidfile:      "/tmp/job.pid",
			MaxProc:      runtime.NumCPU(),
			PprofAddrs:   []string{"localhost:6922"},
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
			{Key: 1, Addr: "tcp@0.0.0.0:6912"},
			{Key: 2, Addr: "tcp@0.0.0.0:6913"},
		},
	}
}
