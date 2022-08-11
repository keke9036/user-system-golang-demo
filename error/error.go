// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/6/13

package errorcode

import (
	"entry-task/util"
	"fmt"
	"net/http"
	"strconv"
)

type Error struct {
	Code int
	Msg  string
}

// 错误码定义
var (
	Success      = NewErrorCode(0, "Success")
	Timeout      = NewErrorCode(1, "ConnectTimeout")
	InternalErr  = NewErrorCode(2, "Internal Server Error")
	ParamInvalid = NewErrorCode(3, "Param invalid")

	UserNameNotValid    = NewErrorCode(101, "用户名不合法")
	PwdNotValid         = NewErrorCode(102, "密码不合法")
	NickNameNotValid    = NewErrorCode(103, "昵称不合法")
	UserLoginErr        = NewErrorCode(104, "用户名或密码错误")
	UserNeedLogin       = NewErrorCode(105, "用户需要登录")
	UserEditErr         = NewErrorCode(106, "用户更新信息失败")
	UserUploadAvatarErr = NewErrorCode(107, "用户上传头像失败")
	UserExist           = NewErrorCode(108, "用户已存在")
	GetUserInfoErr      = NewErrorCode(109, "获取用户信息失败")
	SessionNotExist     = NewErrorCode(110, "用户session不存在")

	PoolClosed       = NewErrorCode(201, "Pool is closed")
	PoolRejected     = NewErrorCode(202, "Pool is rejected")
	RpcClientInitErr = NewErrorCode(203, "Init rpc client fail")
)

var codes = map[string]*Error{}

func NewErrorCode(code int, msg string) *Error {

	codeStr := strconv.Itoa(code)
	if _, ok := codes[codeStr]; ok {
		util.Logger.Errorf("ErrorCode %d exists", code)
	}

	err := Error{Code: code, Msg: msg}
	codes[codeStr] = &err
	return &err
}

func (e *Error) GetError() string {
	return fmt.Sprintf("code: %d, msg: %s", e.GetCode(), e.GetMsg())
}

func (e *Error) GetCode() int {
	return e.Code
}

func (e *Error) GetCodeStr() string {
	return strconv.Itoa(e.Code)
}

func (e *Error) GetMsg() string {
	return e.Msg
}

func GetErrorOrDefault(errMsg string, err *Error) *Error {
	if _, ok := codes[errMsg]; ok {
		return codes[errMsg]
	}

	return err
}

func (e *Error) StatusCode() int {
	switch c := e.GetCode(); {
	case c == Success.GetCode():
		return http.StatusOK
	case c == InternalErr.GetCode():
		return http.StatusInternalServerError
	case c == Timeout.GetCode():
		return http.StatusGatewayTimeout
	case c >= 100:
		return http.StatusOK
	}
	return http.StatusInternalServerError
}
