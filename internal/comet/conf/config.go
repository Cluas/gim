package conf

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/Cluas/gim/pkg/log"
	"github.com/spf13/viper"
)

// Config is struct of comet conf
type Config struct {
	Base            *BaseConfig      `mapstructure:"base"`
	Log             *log.Config      `mapstructure:"log"`
	Websocket       *WebsocketConfig `mapstructure:"websocket"`
	Bucket          *BucketConfig    `mapstructure:"bucket"`
	WriteWait       time.Duration
	PongWait        time.Duration
	PingPeriod      time.Duration
	MaxMessageSize  int64
	ReadBufferSize  int
	WriteBufferSize int
	BroadcastSize   int
}

// BaseConfig is struct of base conf
type BaseConfig struct {
	PidFile         string `mapstructure:"pidfile"`
	ServerId        int8   `mapstructure:"serverId"`
	MaxProc         int
	PprofBind       []string `mapstructure:"pprofBind"` // 性能监控的域名端口
	WriteWait       time.Duration
	PongWait        time.Duration
	PingPeriod      time.Duration
	MaxMessageSize  int64
	BroadcastSize   int
	ReadBufferSize  int
	WriteBufferSize int
	CertPath        string `mapstructure:"certPath"`
	KeyPath         string `mapstructure:"keyPath"`
}
type WebsocketConfig struct {
	Bind string `mapstructure:"bind"` // 性能监控的域名端口
}
type BucketConfig struct {
	Size     int `mapstructure:"size"`
	Channel  int `mapstructure:"channel"`
	Room     int `mapstructure:"room"`
	SvrProto int `mapstructure:"svrProto"`
}

var (
	// Conf is var of conf
	Conf       *Config
	configPath string
)

func init() {
	flag.StringVar(&configPath, "d", "./internal/comet/conf/", "set logic conf file path")
}

// Init is func to initial conf
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
			PidFile:         "/tmp/comet.pid",
			MaxProc:         runtime.NumCPU(),
			WriteWait:       10,
			PongWait:        60,
			PingPeriod:      54,
			MaxMessageSize:  512,
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		Log: &log.Config{
			LogPath:  "./log.log",
			LogLevel: "debug",
		},
		Bucket: &BucketConfig{
			Size:    256,
			Channel: 1024,
			Room:    1024,
		},
		Websocket: &WebsocketConfig{
			Bind: ":7199",
		},
	}
}
