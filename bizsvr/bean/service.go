// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/6/14

package bean

import (
	"context"
	"entry-task/bizsvr/dal"
	"entry-task/bizsvr/service"
	"entry-task/conf"
	"entry-task/util"
)

var (
	UserService *service.UserService
)

func InitUserService(ctx context.Context) {
	dao := dal.NewUserDao(Db, RedisClient)
	idGen := util.NewWorker(conf.RpcServerConfig.WorkerId, conf.RpcServerConfig.WorkerId)

	UserService = service.NewService(ctx, dao, idGen)

}
