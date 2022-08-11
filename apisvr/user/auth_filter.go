// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/7/14

package user

import (
	"context"
	"entry-task/apisvr/base"
	"entry-task/arpc"
	"entry-task/bizsvr/proto"
	"entry-task/conf"
	errorcode "entry-task/error"
	"entry-task/util"
)

var URLS = []string{
	"/api/v1/user/edit",
	"/api/v1/user/uploadAvatar",
	"/api/v1/user/info",
	"/api/v1/user/logout"}

type AuthFilter struct {
	Client *arpc.MClient
}

func NewAuthFilter(client *arpc.MClient) AuthFilter {
	return AuthFilter{
		Client: client,
	}
}

func (f AuthFilter) DoFilter(c *base.GContext) bool {
	util.Logger.Infof("AuthFilter doFilter, url %s", c.URL)

	response := base.NewResp(c)
	sessionId, err := c.Request.Cookie("sessionId")
	if err != nil {
		response.Fail(errorcode.UserNeedLogin)
		return false
	} else {
		var userId string
		ctx, cancel := context.WithTimeout(context.Background(), conf.RpcClientConfig.CallTimeout)
		defer cancel()
		err := f.Client.Call(ctx, proto.SessionGetService, sessionId.Value, &userId)
		if err != nil {
			response.Fail(errorcode.UserNeedLogin)
			return false
		}
		c.Set("userId", userId)
		c.Set("sessionId", sessionId.Value)
	}

	return true
}

func (f AuthFilter) MatchUrl(ctx *base.GContext) bool {

	for _, u := range URLS {
		if ctx.URL == u {
			return true
		}
	}
	return false
}
