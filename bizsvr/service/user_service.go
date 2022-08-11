// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/6/13

package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"entry-task/bizsvr/constant"
	"entry-task/bizsvr/dal"
	"entry-task/bizsvr/entity"
	"entry-task/conf"
	errorcode "entry-task/error"
	"entry-task/util"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	ctx   context.Context
	dao   *dal.UserDao
	idGen *util.Worker
}

func NewService(ctx context.Context, dao *dal.UserDao, idGen *util.Worker) *UserService {
	s := UserService{}
	s.ctx = ctx
	s.dao = dao
	s.idGen = idGen

	return &s
}

func (s UserService) Login(req entity.LoginReq, resp *entity.LoginResp) error {
	util.Logger.WithField("req", req).Infof("User login")
	user, err := s.dao.GetUserByName(req.UserName)
	if err != nil {
		util.Logger.Errorf("GetUserByName error, userName %s, error %v", req.UserName, err)
		return err
	}

	// check password match
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		util.Logger.Errorf("password not match, userName %s", req.UserName)
		return err
	}

	sessionId := util.GenerateSid()
	resp.SessionId = sessionId
	userId := util.IntToStr(user.UserId)

	err = s.dao.SetCache(s.ctx, constant.SessionIdCachePrefix+sessionId, userId, constant.ThreeDaysToExpire)
	if err != nil {
		util.Logger.Errorf("set cache erro, userName %s, error %v", req.UserName, err)
		return err
	}
	return nil
}

func (s UserService) GetUser(req entity.InfoReq, resp *entity.UserInfoResp) error {
	util.Logger.WithField("req", req).Infof("User Get")
	var userId string
	err := s.GetSession(req.SessionId, &userId)
	if err != nil {
		util.Logger.Errorf("GetCache user info from cache error, userId %s, error %v", userId, err)
		return err
	}

	userCacheKey := constant.UserInfoCachePrefix + userId
	userCache, err := s.dao.GetCache(s.ctx, userCacheKey)
	if err == nil {
		err = json.Unmarshal([]byte(userCache), &resp)
		if err != nil {
			util.Logger.Errorf("User info parse json error, userId %s, %v", userId, err)
		} else {
			return nil
		}
	}
	uid, err := util.StrToInt(userId)
	if err != nil {
		util.Logger.Errorf("GetCache userId invalid, userId %s, error %v", userId, err)
		return err
	}
	user, err := s.dao.GetUserByUserId(uid)
	if err != nil {
		util.Logger.Errorf("GetCache user info from db error, userId %s, error %v", userId, err)
		return err
	}
	resp.UserName = user.UserName
	resp.NickName = user.NickName
	if user.AvatarUrl.Valid {
		resp.AvatarUrl = user.AvatarUrl.String
	} else {
		resp.AvatarUrl = ""
	}
	resp.UserId = userId

	userToCache, _ := json.Marshal(resp)
	err = s.dao.SetCache(s.ctx, userCacheKey, string(userToCache), constant.ThreeDaysToExpire)
	if err != nil {
		util.Logger.Errorf("cache user fail, userId %s, error %v", userId, err)
	}
	return nil
}

func (s UserService) EditUser(req entity.EditReq, resp *entity.EditResp) error {
	util.Logger.WithField("req", req).Infof("User Edit")
	var userId string
	err := s.GetSession(req.SessionId, &userId)
	if err != nil {
		util.Logger.Errorf("GetCache user info from cache error, userId %s, error %v", userId, err)
		return err
	}

	uid, err := util.StrToInt(userId)
	if err != nil {
		util.Logger.Errorf("GetCache userId invalid, userId %s, error %v", userId, err)
		return err
	}

	_, err = s.dao.GetUserByUserId(uid)
	if err != nil {
		return errors.New(errorcode.GetUserInfoErr.GetCodeStr())
	}

	userToUpdate := dal.UserDo{
		UserId:   uid,
		NickName: req.NickName,
	}
	if len(req.AvatarUrl) > 0 {
		userToUpdate.AvatarUrl = sql.NullString{String: req.AvatarUrl, Valid: true}
	} else {
		userToUpdate.AvatarUrl = sql.NullString{Valid: false}
	}

	err = s.dao.UpdateUser(userToUpdate)
	if err != nil {
		return err
	}

	cacheKey := constant.UserInfoCachePrefix + userId
	err = s.dao.DeleteCache(s.ctx, cacheKey)
	if err != nil {
		util.Logger.Errorf("DeleteCache user cache error, userId %s, error %v", userId, err)
	}

	return nil
}

func (s UserService) UploadAvatar(req entity.UploadReq, resp *entity.UploadResp) error {
	util.Logger.WithField("req", req).Infof("User Upload Avatar")
	fileName := util.BuildFileName(req.FileName)
	dir := constant.UploadFileDir + "/" + req.UserName + "/"
	dst := dir + "/" + fileName

	err := util.SaveFile(&req.Content, dst, dir)
	if err != nil {
		return err
	}

	resp.AvatarUrl = conf.WebConf.StaticFileUrl + "/" + req.UserName + "/" + fileName
	return nil
}

func (s UserService) Register(req entity.RegisterReq, resp *entity.RegisterResp) error {
	util.Logger.WithField("req", req).Infof("User register")
	_, err := s.dao.GetUserByName(req.UserName)
	if !errors.Is(err, sql.ErrNoRows) {
		return errors.New(errorcode.UserExist.GetMsg())
	}

	var u dal.UserDo
	u.UserName = req.UserName
	u.Password = req.Password
	u.NickName = req.NickName
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 4)
	if err != nil {
		return err
	}
	u.Password = string(bytes)

	// generate user id
	id, err := s.idGen.NextID()
	if err != nil {
		util.Logger.Errorf("Generate id error, %v", err)
		return errors.New(errorcode.InternalErr.GetMsg())
	}
	u.UserId = id

	err = s.dao.AddUser(u)
	if err != nil {
		return err
	}
	return nil
}

func (s UserService) GetSession(sessionId string, userId *string) error {
	util.Logger.WithField("sessionId", sessionId).
		WithField("userId", userId).
		Infof("User get session")
	key := constant.SessionIdCachePrefix + sessionId
	var err error
	*userId, err = s.dao.GetCache(s.ctx, key)
	if err != nil {
		return errors.New(errorcode.SessionNotExist.GetCodeStr())
	}

	return nil
}

func (s UserService) Logout(req entity.LogoutReq, resp *entity.LogoutResp) error {
	util.Logger.WithField("req", req).Infof("User logout")
	var userId string
	err := s.GetSession(req.SessionId, &userId)
	if err != nil {
		util.Logger.Errorf("GetCache user info from cache error, SessionId %s, error %v", req.SessionId, err)
		return err
	}

	// delete session
	cacheKey := constant.SessionIdCachePrefix + req.SessionId
	err = s.dao.DeleteCache(s.ctx, cacheKey)
	if err != nil {
		util.Logger.Errorf("DeleteCache user cache error, sessionId %s, error %v", req.SessionId, err)
		return err
	}

	return nil
}
