package conf

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/viper"

	"github.com/Cluas/gim/pkg/log"
)

// Config is struct of comet conf
type Config struct {
	Base      *BaseConf      `mapstructure:"base"`
	Log       *log.Config    `mapstructure:"log"`
	Websocket *WebsocketConf `mapstructure:"websocket"`
	Bucket    *BucketConf    `mapstructure:"bucket"`
	RPC       *RPCConf       `mapstructure:"rpc"`
}

// BaseConf is struct of base conf
type BaseConf struct {
	PidFile   string   `mapstructure:"pid_file"`
	ServerID  string   `mapstructure:"server_id"`
	PprofBind []string `mapstructure:"pprof_bind"`
	CertPath  string   `mapstructure:"cert_path"`
	KeyPath   string   `mapstructure:"key_path"`
	MaxProc   int      `mapstructure:"max_process"`
}

// WebsocketConf is struct of websocket conf
type WebsocketConf struct {
	Bind            string        `mapstructure:"port"`
	WriteWait       time.Duration `mapstructure:"write_wait"`
	PongWait        time.Duration `mapstructure:"pong_wait"`
	PingPeriod      time.Duration `mapstructure:"ping_period"`
	MaxMessageSize  int64         `mapstructure:"max_message_size"`
	ReadBufferSize  int           `mapstructure:"read_buffer_size"`
	WriteBufferSize int           `mapstructure:"write_buffer_size"`
}

// BucketConf is struct of BucketConf
type BucketConf struct {
	Size          int    `mapstructure:"size"`
	Channel       int    `mapstructure:"channel"`
	Room          int    `mapstructure:"room"`
	RoutineAmount uint64 `mapstructure:"routine_amount"`
	RoutineSize   int    `mapstructure:"routine_size"`
	BroadcastSize int    `mapstructure:"broadcast_size"`
}

//RPCConf is struct of RPCConf
type RPCConf struct {
	LogicAddr []Address `mapstructure:"logic_bind"`
	CometAddr []Address `mapstructure:"comet_bind"`
}

// Address is struct of rpc address
type Address struct {
	Key  int8   `mapstructure:"key"`
	Addr string `mapstructure:"addr"`
}

var (
	// Conf is var of conf
	Conf       *Config
	configFile string
)

func init() {
	flag.StringVar(&configFile, "c", ".", "set logic conf file.")
}

// Init is func to initial conf
func Init() (err error) {
	Conf = NewConfig()
	viper.SetConfigFile(configFile)

	if err = viper.ReadInConfig(); err != nil {
		return err
	}

	if err = viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Errorf("unable to decode into struct:  %s \n ", err))
	}
	return nil
}

// NewConfig is constructor of Config
func NewConfig() *Config {
	return &Config{
		Base: &BaseConf{
			PidFile:  "/tmp/comet.pid",
			MaxProc:  runtime.NumCPU(),
			ServerID: "1",
		},
		Log: &log.Config{
			ServiceName: "comet",
		},
		Bucket: &BucketConf{
			Size:          8,
			Channel:       1024,
			Room:          1024,
			RoutineAmount: 32,
			RoutineSize:   20,
		},
		RPC: &RPCConf{
			LogicAddr: []Address{{Addr: "tcp@0.0.0.0:6923", Key: 1}},
			CometAddr: []Address{{Addr: "tcp@0.0.0.0:6999", Key: 1}},
		},
		Websocket: &WebsocketConf{
			Bind:            ":7199",
			MaxMessageSize:  512,
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			PingPeriod:      54 * time.Second,
			PongWait:        60 * time.Second,
			WriteWait:       10 * time.Second,
		},
	}
}
