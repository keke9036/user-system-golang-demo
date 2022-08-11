package conf

import (
	"github.com/spf13/viper"
	"time"
)

var (
	WebConf         *WebServerConf
	DbConf          *DbServerConf
	CacheConf       *CacheServerConf
	RpcServerConfig *RpcServerConf
	RpcClientConfig *RpcClientConf
)

const confExtension = "yaml"
const confDir = "conf"

func LoadHttpServerConf(confFile string) error {
	vp := viper.New()
	vp.SetConfigName(confFile)
	vp.SetConfigType(confExtension)
	vp.AddConfigPath(confDir)

	err := vp.ReadInConfig()
	if err != nil {
		return err
	}

	// web服务配置
	err = vp.UnmarshalKey("WebServer", &WebConf)
	if err != nil {
		return err
	}

	WebConf.ReadTimeout *= time.Second
	WebConf.WriteTimeout *= time.Second

	// Rpc
	err = vp.UnmarshalKey("RpcClient", &RpcClientConfig)
	if err != nil {
		return err
	}
	RpcClientConfig.CallTimeout *= time.Millisecond
	RpcClientConfig.ConnectTimeout *= time.Millisecond

	return nil
}

func LoadRpcServerConf(confFile string) error {
	vp := viper.New()
	vp.SetConfigName(confFile)
	vp.SetConfigType(confExtension)
	vp.AddConfigPath(confDir)

	err := vp.ReadInConfig()
	if err != nil {
		return err
	}

	// DB配置
	err = vp.UnmarshalKey("DbServer", &DbConf)
	if err != nil {
		return err
	}

	// Redis配置
	err = vp.UnmarshalKey("CacheServer", &CacheConf)
	if err != nil {
		return err
	}
	CacheConf.IdleTimeout *= time.Second

	// Rpc
	err = vp.UnmarshalKey("RpcServer", &RpcServerConfig)
	if err != nil {
		return err
	}

	RpcServerConfig.ServerHandleTimeout *= time.Millisecond

	return nil
}
