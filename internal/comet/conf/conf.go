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
	Base            *BaseConf      `mapstructure:"base"`
	Log             *log.Config    `mapstructure:"log"`
	Websocket       *WebsocketConf `mapstructure:"websocket"`
	Bucket          *BucketConf    `mapstructure:"bucket"`
	WriteWait       time.Duration
	PongWait        time.Duration
	PingPeriod      time.Duration
	MaxMessageSize  int64
	ReadBufferSize  int
	WriteBufferSize int
	BroadcastSize   int
	RPC             *RPCConf
}

// BaseConf is struct of base conf
type BaseConf struct {
	PidFile    string `mapstructure:"pidfile"`
	ServerID   int8   `mapstructure:"serverID"`
	MaxProc    int
	PprofBind  []string `mapstructure:"pprofBind"` // 性能监控的域名端口
	WriteWait  time.Duration
	PongWait   time.Duration
	PingPeriod time.Duration
	CertPath   string `mapstructure:"certPath"`
	KeyPath    string `mapstructure:"keyPath"`
}

// WebsocketConf is struct of websocket conf
type WebsocketConf struct {
	Bind string `mapstructure:"bind"` // 性能监控的域名端口
}

// BucketConf is struct of BucketConf
type BucketConf struct {
	Size     int `mapstructure:"size"`
	Channel  int `mapstructure:"channel"`
	Room     int `mapstructure:"room"`
	SvrProto int `mapstructure:"svrProto"`
}

//RPCConf is struct of RPCConf
type RPCConf struct {
	LogicAddr []Address `mapstructure:"rpcLogicAddrs"`
	CometAddr []Address `mapstructure:"comet_addr"`
}

// Address is struct of rpc address
type Address struct {
	Key  int8   `mapstructure:"key"`
	Addr string `mapstructure:"addr"`
}

var (
	// Conf is var of conf
	Conf       *Config
	configPath string
)

func init() {
	flag.StringVar(&configPath, "p", ".", "set logic conf file path")
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
		Base: &BaseConf{
			PidFile:    "/tmp/comet.pid",
			MaxProc:    runtime.NumCPU(),
			WriteWait:  10,
			PongWait:   60,
			PingPeriod: 54,
		},
		Log: &log.Config{
			LogPath:  "./log.log",
			LogLevel: "debug",
		},
		Bucket: &BucketConf{
			Size:    256,
			Channel: 1024,
			Room:    1024,
		},
		RPC: &RPCConf{
			LogicAddr: []Address{{Addr: "tcp@0.0.0.0:6923", Key: 1}},
			CometAddr: []Address{{Addr: "tcp@0.0.0.0:6912", Key: 1}},
		},
		Websocket: &WebsocketConf{
			Bind: ":7199",
		},
		MaxMessageSize:  512,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		PingPeriod:      54 * time.Second,
		PongWait:        60 * time.Second,
		WriteWait:       10 * time.Second,
	}
}
