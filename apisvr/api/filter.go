// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/6/13

package api

import (
	"context"
	"entry-task/bizsvr/proto"
	"entry-task/conf"
	"entry-task/error"
	"github.com/gin-gonic/gin"
)

func SessionRequired(c *gin.Context) {
	response := NewResp(c)
	sessionId, err := c.Cookie("sessionId")
	if err != nil {
		response.FailAndAbort(errorcode.UserNeedLogin)
		return
	} else {
		var userId string
		ctx, cancel := context.WithTimeout(context.Background(), conf.RpcClientConfig.CallTimeout)
		defer cancel()
		err := RpcClient.Call(ctx, proto.SessionGetService, sessionId, &userId)
		if err != nil {
			response.FailAndAbort(errorcode.UserNeedLogin)
			return
		}
		c.Set("userId", userId)
		c.Set("sessionId", sessionId)
		c.Next()
	}
}
