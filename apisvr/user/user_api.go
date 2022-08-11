// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/7/18

package user

import (
	"context"
	"encoding/json"
	"entry-task/apisvr/base"
	"entry-task/arpc"
	"entry-task/bizsvr/constant"
	"entry-task/bizsvr/entity"
	"entry-task/bizsvr/proto"
	"entry-task/conf"
	errorcode "entry-task/error"
	"entry-task/util"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

type User struct {
	Client *arpc.MClient
}

type Void struct{}

func NewUser(client *arpc.MClient) User {
	return User{Client: client}
}

/*************** 自定义handler ***************/

// LoginHandler 用户登录
func (u User) LoginHandler(c *base.GContext) {

	req := entity.LoginReq{}
	resp := base.NewResp(c)
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		resp.Fail(errorcode.ParamInvalid)
		return
	}
	defer c.Request.Body.Close()

	err = json.Unmarshal(body, &req)
	if err != nil {
		resp.Fail(errorcode.ParamInvalid)
		return
	}

	util.Logger.Infof("user %s login", req.UserName)
	data := entity.LoginResp{}

	ctx, cancel := context.WithTimeout(context.Background(), conf.RpcClientConfig.CallTimeout)
	defer cancel()
	err = u.Client.Call(ctx, proto.UserLoginService, req, &data)
	if err != nil {
		resp.Fail(errorcode.UserLoginErr)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "sessionId",
		Value:    data.SessionId,
		Expires:  time.Now().Add(time.Hour),
		Path:     "/",
		Domain:   "",
		Secure:   false,
		HttpOnly: true,
	})
	resp.Success(data)
}

func (u User) RegisterHandler(c *base.GContext) {
	req := entity.RegisterReq{}
	resp := base.NewResp(c)

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		resp.Fail(errorcode.ParamInvalid)
		return
	}

	defer c.Request.Body.Close()
	err = json.Unmarshal(body, &req)
	if err != nil {
		util.Logger.Errorf("user register param parse error, %v", err)
		resp.Fail(errorcode.ParamInvalid)
		return
	}

	data := entity.RegisterResp{}
	ctx, cancel := context.WithTimeout(context.Background(), conf.RpcClientConfig.CallTimeout)
	defer cancel()
	err = u.Client.Call(ctx, proto.UserRegisterService, req, &data)
	if err != nil {
		util.Logger.Errorf("user register error, %v", err)
		resp.FailWithMsg(err.Error())
		return
	}

	resp.Success(data)
}

// InfoHandler 获取用户信息
func (u User) InfoHandler(c *base.GContext) {
	resp := base.NewResp(c)

	sessionId, ok := c.Get("sessionId")
	if !ok {
		resp.Fail(errorcode.SessionNotExist)
		return
	}
	util.Logger.Infof("GetCache user info, sessionId %s", sessionId)

	req := entity.InfoReq{}
	req.SessionId = fmt.Sprintf("%v", sessionId)
	data := entity.UserInfoResp{}

	ctx, cancel := context.WithTimeout(context.Background(), conf.RpcClientConfig.CallTimeout)
	defer cancel()
	err := u.Client.Call(ctx, proto.UserGetService, req, &data)
	if err != nil {
		util.Logger.Errorf("GetCache user info error, sessionId %s, %v", sessionId, err)
		resp.Fail(errorcode.GetErrorOrDefault(err.Error(), errorcode.GetUserInfoErr))
		return
	}

	resp.Success(data)
}

// EditHandler 编辑用户信息
func (u User) EditHandler(c *base.GContext) {
	req := entity.EditReq{}
	resp := base.NewResp(c)

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		resp.Fail(errorcode.ParamInvalid)
		return
	}

	defer c.Request.Body.Close()
	err = json.Unmarshal(body, &req)
	if err != nil {
		util.Logger.Errorf("user edit param parse error, %v", err)
		resp.Fail(errorcode.ParamInvalid)
		return
	}

	sessionId, _ := c.Get("sessionId")
	req.SessionId = fmt.Sprintf("%v", sessionId)
	util.Logger.Infof("Edit user info, sessionId %s, nickName %s, avatarUrl %s",
		sessionId, req.NickName, req.AvatarUrl)

	data := entity.EditResp{}
	ctx, cancel := context.WithTimeout(context.Background(), conf.RpcClientConfig.CallTimeout)
	defer cancel()
	err = u.Client.Call(ctx, proto.UserEditService, req, &data)
	if err != nil {
		util.Logger.Errorf("Edit user info error, sessionId %s, %v", sessionId, err)
		resp.Fail(errorcode.GetErrorOrDefault(err.Error(), errorcode.UserEditErr))
		return
	}

	resp.Success(Void{})
}

// UploadAvatarHandler 上传头像
func (u User) UploadAvatarHandler(c *base.GContext) {
	resp := base.NewResp(c)

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		util.Logger.Errorf("Upload avatar error, %v", err)
		resp.Fail(errorcode.UserUploadAvatarErr)
		return
	}

	if fileHeader == nil {
		util.Logger.Errorf("Upload avatar fileHeader null, %v", err)
		resp.Fail(errorcode.UserUploadAvatarErr)
		return
	}

	// Basic Check
	fileExt := util.GetFileExt(fileHeader.Filename)
	fileName := fileHeader.Filename
	if !checkExt(fileExt) {
		util.Logger.Errorf("Upload avatar file extensioin not support, %s", fileExt)
		resp.Fail(errorcode.UserUploadAvatarErr)
		return
	}

	if !checkFileSize(file) {
		util.Logger.Errorf("Upload avatar file size greater than %d byte", constant.MaxImageSize)
		resp.Fail(errorcode.UserUploadAvatarErr)
		return
	}

	content, err := util.ReadFileByte(fileHeader)
	if err != nil {
		util.Logger.Error("Upload avatar file read error")
		resp.Fail(errorcode.UserUploadAvatarErr)
		return
	}

	userId, _ := c.Get("userId")
	req := entity.UploadReq{
		UserName: fmt.Sprintf("%v", userId),
		FileName: fileName,
		FileExt:  fileExt,
		Content:  content,
	}
	data := entity.UploadResp{}
	ctx, cancel := context.WithTimeout(context.Background(), conf.RpcClientConfig.CallTimeout)
	defer cancel()
	err = u.Client.Call(ctx, proto.AvatarUploadService, req, &data)
	if err != nil {
		resp.Fail(errorcode.UserUploadAvatarErr)
		return
	}

	resp.Success(data)
}

func (u User) LogoutHandler(c *base.GContext) {
	resp := base.NewResp(c)

	sessionId, _ := c.Get("sessionId")
	util.Logger.Infof("GetCache user info, sessionId %s", sessionId)

	req := entity.LogoutReq{}
	req.SessionId = fmt.Sprintf("%v", sessionId)
	data := entity.LogoutResp{}

	ctx, cancel := context.WithTimeout(context.Background(), conf.RpcClientConfig.CallTimeout)
	defer cancel()
	err := u.Client.Call(ctx, proto.UserLogoutService, req, &data)
	if err != nil {
		util.Logger.Errorf("Logout error, sessionId %s, %v", sessionId, err)
		resp.Fail(errorcode.GetErrorOrDefault(err.Error(), errorcode.GetUserInfoErr))
		return
	}

	resp.Success(data)
}

func (u User) PingHandler(c *base.GContext) {
	resp := base.NewResp(c)
	data := entity.PingResp{
		Text: "hello",
	}
	resp.Success(data)
}

func checkExt(ext string) bool {
	for _, allowExt := range constant.ImageExtSupports {
		if strings.EqualFold(allowExt, ext) {
			return true
		}
	}
	return false
}

func checkFileSize(f multipart.File) bool {
	content, _ := ioutil.ReadAll(f)
	size := len(content)

	return size <= constant.MaxImageSize
}
