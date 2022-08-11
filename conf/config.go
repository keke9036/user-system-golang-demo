package conf

import "time"

/**
 * apisvr server配置
 */
type WebServerConf struct {
	Addr          string
	Port          string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	StaticFileUrl string
}

/**
 * MySQL配置
 */
type DbServerConf struct {
	Username string
	Password string
	Host     string
	Port     string
	DbName   string
}

/**
 * Redis配置
 */
type CacheServerConf struct {
	Host        string
	DbIndex     int
	PoolSize    int
	MaxRetries  int
	IdleTimeout time.Duration
}

/**
 * RPC配置
 */
type RpcServerConf struct {
	Addr                string
	ServerHandleTimeout time.Duration
	WorkerId            int64
}

type RpcClientConf struct {
	Addrs          []string      // rpc server列表
	CallTimeout    time.Duration // 接口调用超时
	ConnectTimeout time.Duration
	MaxConnections int
}
