package main

import (
	"entry-task/apisvr/user"
	"entry-task/conf"
	"entry-task/util"
	"flag"
	huge "github.com/dablelv/go-huge-util"
	"github.com/sirupsen/logrus"
	"net/http"
)

/*
 基于gin框架实现http功能
*/

func main() {

	var configFile = flag.String("config", "config_httpserver", "config file")
	var logLevel = flag.String("logLevel", "Info", "log level")
	var logPath = flag.String("logPath", "./logs/api-server.log", "log path")
	flag.Parse()

	// 初始化Log
	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		util.Logger.Warnf("Log level illegal, set to info")
		level = logrus.InfoLevel
	}
	util.InitLog(*logPath, level)
	util.Logger.Info("App bean start")

	// 初始化配置信息
	util.Logger.Info("Init config")

	err = conf.LoadHttpServerConf(*configFile)
	if err != nil {
		util.Logger.Fatalf("Read conf error: %v", err)
		return
	}
	printConf()

	util.Logger.Info("Start biz api server")
	err = startWebServer()
	if err != nil {
		util.Logger.Fatalf("Start biz api server error: %v", err)
		return
	}
}

/**
 * 启动web server
 */
func startWebServer() error {
	//gin.SetMode(gin.ReleaseMode)
	//r := api.InitRouter()

	// 自定义dispatcher
	handler, err := user.InitHandler()
	if err != nil {
		util.Logger.Errorf("InitHandler error, %v", err)
		return err
	}

	s := &http.Server{
		Addr:    conf.WebConf.Addr + ":" + conf.WebConf.Port,
		Handler: handler,
		//Handler:      initHandler(),
		ReadTimeout:  conf.WebConf.ReadTimeout,
		WriteTimeout: conf.WebConf.WriteTimeout,
	}

	err = s.ListenAndServe()
	if err != nil {
		return err
	}

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
