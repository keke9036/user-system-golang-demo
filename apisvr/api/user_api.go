package api

import (
	"context"
	"entry-task/bizsvr/constant"
	"entry-task/bizsvr/entity"
	"entry-task/bizsvr/proto"
	"entry-task/conf"
	"entry-task/error"
	"entry-task/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"strings"
)

type User struct {
}

type Void struct{}

func NewUser() User {
	return User{}
}

// Login 用户登录
func (u User) Login(c *gin.Context) {

	req := entity.LoginReq{}
	resp := NewResp(c)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		util.Logger.Errorf("user login param parse error, %v", err)
		resp.Fail(errorcode.ParamInvalid)
		return
	}

	util.Logger.Infof("user %s login", req.UserName)
	data := entity.LoginResp{}

	ctx, cancel := context.WithTimeout(context.Background(), conf.RpcClientConfig.CallTimeout)
	defer cancel()
	err = RpcClient.Call(ctx, proto.UserLoginService, req, &data)
	if err != nil {
		resp.Fail(errorcode.UserLoginErr)
		return
	}

	c.SetCookie("sessionId", data.SessionId, 3600, "/", "", false, true)
	resp.Success(data)
}

// Info 获取用户信息
func (u User) Info(c *gin.Context) {
	resp := NewResp(c)

	sessionId, _ := c.Get("sessionId")
	util.Logger.Infof("GetCache user info, sessionId %s", sessionId)

	req := entity.InfoReq{}
	req.SessionId = fmt.Sprintf("%v", sessionId)
	data := entity.UserInfoResp{}

	ctx, cancel := context.WithTimeout(context.Background(), conf.RpcClientConfig.CallTimeout)
	defer cancel()
	err := RpcClient.Call(ctx, proto.UserGetService, req, &data)
	if err != nil {
		util.Logger.Errorf("GetCache user info error, sessionId %s, %v", sessionId, err)
		resp.Fail(errorcode.GetErrorOrDefault(err.Error(), errorcode.GetUserInfoErr))
		return
	}

	resp.Success(data)
}

// Edit 编辑用户信息
func (u User) Edit(c *gin.Context) {
	req := entity.EditReq{}
	resp := NewResp(c)
	err := c.ShouldBindJSON(&req)
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
	err = RpcClient.Call(ctx, proto.UserEditService, req, &data)
	if err != nil {
		util.Logger.Errorf("Edit user info error, sessionId %s, %v", sessionId, err)
		resp.Fail(errorcode.GetErrorOrDefault(err.Error(), errorcode.UserEditErr))
		return
	}

	resp.Success(Void{})
}

// UploadAvatar 上传头像
func (u User) UploadAvatar(c *gin.Context) {
	resp := NewResp(c)

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		util.Logger.Errorf("Upload avatar error, %v", err)
		resp.Fail(errorcode.UserUploadAvatarErr)
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
	err = RpcClient.Call(ctx, proto.AvatarUploadService, req, &data)
	if err != nil {
		resp.Fail(errorcode.UserUploadAvatarErr)
		return
	}

	resp.Success(data)
}

// Register 注册用户，测试用
func (u User) Register(c *gin.Context) {
	req := entity.RegisterReq{}
	resp := NewResp(c)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		util.Logger.Errorf("user register param parse error, %v", err)
		resp.Fail(errorcode.ParamInvalid)
		return
	}

	data := entity.RegisterResp{}
	ctx, cancel := context.WithTimeout(context.Background(), conf.RpcClientConfig.CallTimeout)
	defer cancel()
	err = RpcClient.Call(ctx, proto.UserRegisterService, req, &data)
	if err != nil {
		util.Logger.Errorf("user register error, %v", err)
		resp.FailWithMsg(err.Error())
		return
	}

	resp.Success(data)
}

func (u User) Logout(c *gin.Context) {
	resp := NewResp(c)

	sessionId, _ := c.Get("sessionId")
	util.Logger.Infof("GetCache user info, sessionId %s", sessionId)

	req := entity.LogoutReq{}
	req.SessionId = fmt.Sprintf("%v", sessionId)
	data := entity.LogoutResp{}

	ctx, cancel := context.WithTimeout(context.Background(), conf.RpcClientConfig.CallTimeout)
	defer cancel()
	err := RpcClient.Call(ctx, proto.UserLogoutService, req, &data)
	if err != nil {
		util.Logger.Errorf("Logout error, sessionId %s, %v", sessionId, err)
		resp.Fail(errorcode.GetErrorOrDefault(err.Error(), errorcode.GetUserInfoErr))
		return
	}

	resp.Success(data)
}

func (u User) Ping(c *gin.Context) {
	resp := NewResp(c)
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
