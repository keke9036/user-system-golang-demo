// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/6/13

package api

import (
	"entry-task/error"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Resp struct {
	Ctx *gin.Context
}

func NewResp(ctx *gin.Context) *Resp {
	return &Resp{
		Ctx: ctx,
	}
}

func (r *Resp) Success(data interface{}) {
	resp := gin.H{"code": 0, "data": data}
	r.Ctx.JSONP(http.StatusOK, resp)
}

func (r *Resp) Fail(err *errorcode.Error) {
	response := gin.H{"code": err.GetCode(), "msg": err.GetMsg()}
	r.Ctx.JSON(err.StatusCode(), response)
}

func (r *Resp) FailWithMsg(msg string) {
	response := gin.H{"code": errorcode.InternalErr.GetCode(), "msg": msg}
	r.Ctx.JSON(http.StatusInternalServerError, response)
}

func (r *Resp) FailAndAbort(err *errorcode.Error) {
	response := gin.H{"code": err.GetCode(), "msg": err.GetMsg()}
	r.Ctx.AbortWithStatusJSON(err.StatusCode(), response)
}
