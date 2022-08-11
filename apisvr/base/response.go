// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/6/13

package base

import (
	"encoding/json"
	"entry-task/error"
	"net/http"
)

type Resp struct {
	Ctx *GContext
}

type H map[string]any

func NewResp(ctx *GContext) *Resp {
	return &Resp{
		Ctx: ctx,
	}
}

func (r *Resp) Success(data interface{}) {
	r.Ctx.Writer.Header().Set("Content-Type", "application/json")
	resp := H{"code": 0, "data": data}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		r.Ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	r.Ctx.Writer.Write(jsonResp)
}

func (r *Resp) Fail(error *errorcode.Error) {
	r.Ctx.Writer.Header().Set("Content-Type", "application/json")

	resp := H{"code": error.GetCode(), "msg": error.GetMsg()}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		r.Ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		r.Ctx.Writer.WriteHeader(error.StatusCode())
	}
	r.Ctx.Writer.Write(jsonResp)
}

func (r *Resp) FailWithMsg(msg string) {
	r.Ctx.Writer.Header().Set("Content-Type", "application/json")
	r.Ctx.Writer.WriteHeader(http.StatusInternalServerError)

	resp := H{"code": errorcode.InternalErr.GetCode(), "msg": msg}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		r.Ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	r.Ctx.Writer.Write(jsonResp)

}
