// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/7/18

package user

import (
	"entry-task/apisvr/base"
	"entry-task/arpc"
	"entry-task/bizsvr/constant"
	"entry-task/conf"
	errorcode "entry-task/error"
	"entry-task/util"
	"errors"
)

func InitHandler() (*base.Dispatcher, error) {
	util.Logger.Info("Init rpc client")

	client := arpc.NewMClientAndInit(conf.RpcClientConfig.Addrs, arpc.RandomSelect, conf.RpcClientConfig.ConnectTimeout)
	if client == nil {
		util.Logger.WithField("addrs", conf.RpcClientConfig.Addrs).
			WithField("timeout", conf.RpcClientConfig.ConnectTimeout).
			Errorf("NewMClientAndInit fail")
		return nil, errors.New(errorcode.RpcClientInitErr.GetMsg())
	}

	dispatcher := base.NewDispatcher()
	dispatcher.RegisterPageDir("/static", "./static")
	dispatcher.RegisterPageDir("/upload", constant.UploadFileDir)
	dispatcher.RegisterPageFile("/index", "./static/html/login.html")

	dispatcher.RegisterFilter(NewAuthFilter(client))

	u := NewUser(client)
	dispatcher.RegisterHandler(userPrefix("/login"), base.Post, u.LoginHandler)
	dispatcher.RegisterHandler(userPrefix("/register"), base.Post, u.RegisterHandler)
	dispatcher.RegisterHandler(userPrefix("/ping"), base.Get, u.PingHandler)

	dispatcher.RegisterHandler(userPrefix("/edit"), base.Post, u.EditHandler)
	dispatcher.RegisterHandler(userPrefix("/uploadAvatar"), base.Post, u.UploadAvatarHandler)
	dispatcher.RegisterHandler(userPrefix("/logout"), base.Post, u.LogoutHandler)

	dispatcher.RegisterHandler("/api/v1/user/info", base.Get, u.InfoHandler)

	return dispatcher, nil
}

func userPrefix(url string) string {
	return "/api/v1/user" + url
}
