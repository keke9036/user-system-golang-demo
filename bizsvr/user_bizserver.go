package main

import (
	"context"
	"entry-task/arpc"
	"entry-task/bizsvr/bean"
	"entry-task/bizsvr/cron"
	"entry-task/conf"
	"entry-task/util"
	"flag"
	huge "github.com/dablelv/go-huge-util"
	"github.com/sirupsen/logrus"
	"net"
)

func main() {

	var configFile = flag.String("config", "config_rpcserver.yaml", "config file")
	var logLevel = flag.String("logLevel", "Info", "log level")
	var logPath = flag.String("logPath", "stdout", "log path")

	flag.Parse()

	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		util.Logger.Warnf("Log level illegal, set to info")
		level = logrus.InfoLevel
	}

	// 初始化Log
	util.InitLog(*logPath, level)
	util.Logger.Info("App bean start")

	// 初始化配置信息
	util.Logger.Info("Init config")
	err = conf.LoadRpcServerConf(*configFile)
	if err != nil {
		util.Logger.Fatalf("Read conf error: %v", err)
		return
	}
	printConf()

	util.Logger.Info("Init db")
	err = bean.InitDb(conf.DbConf)
	if err != nil {
		util.Logger.Fatalf("Init db error: %v", err)
		return
	}

	util.Logger.Info("Init cache")
	err = bean.InitCache(conf.CacheConf)
	if err != nil {
		util.Logger.Fatalf("Init cache error: %v", err)
		return
	}

	// 定时清理用户上传图片任务
	cron.RunCleanTask()

	util.Logger.Info("Init rpc service")
	ctx := context.Background()
	bean.InitUserService(ctx)
	err = registerService()
	if err != nil {
		util.Logger.Fatal("register service error:", err)
		return
	}

	// start server
	err = startRpcServer(conf.RpcServerConfig)
	if err != nil {
		util.Logger.Fatal("start rpc server error:", err)
		return
	}
}

func registerService() error {
	if err := arpc.Register(bean.UserService); err != nil {
		return err
	}
	return nil
}

func startRpcServer(serverConf *conf.RpcServerConf) error {

	l, err := net.Listen("tcp", serverConf.Addr)
	if err != nil {
		util.Logger.Fatal("network error:", err)
		return err
	}
	util.Logger.Println("start arpc server on ", l.Addr())
	arpc.Accept(l, serverConf.ServerHandleTimeout)
	return nil
}

func printConf() {
	webConfStr, _ := huge.ToIndentJSON(&conf.WebConf)
	util.Logger.Infof("Web conf: %s", webConfStr)
	dbConfStr, _ := huge.ToIndentJSON(&conf.DbConf)
	util.Logger.Infof("MySQL conf: %s", dbConfStr)
	cacheConfStr, _ := huge.ToIndentJSON(&conf.CacheConf)
	util.Logger.Infof("Redis conf: %s", cacheConfStr)
	rpcConfStr, _ := huge.ToIndentJSON(&conf.RpcServerConfig)
	util.Logger.Infof("Rpc conf: %s", rpcConfStr)
}
