// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/6/17

package service

import (
	"context"
	"encoding/json"
	"entry-task/bizsvr/constant"
	"entry-task/bizsvr/dal"
	"entry-task/bizsvr/entity"
	"entry-task/conf"
	errorcode "entry-task/error"
	"entry-task/util"
	"errors"
	"fmt"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/magiconair/properties/assert"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"testing"
	"time"
)

func TestUserService_Login(t *testing.T) {
	s := NewService(context.Background(), nil, nil)

	userName := "test"
	password := "rootpwd"

	sessionId := "test_sessionId"
	req := entity.LoginReq{
		UserName: userName,
		Password: password,
	}
	resp := entity.LoginResp{}

	t.Run("Login success", func(t *testing.T) {
		// mock dao
		patches := gomonkey.ApplyMethod(reflect.TypeOf(s.dao),
			"GetUserByName",
			func(_ *dal.UserDao, _ string) (dal.UserDo, error) {
				bytes, err := bcrypt.GenerateFromPassword([]byte(password), 6)
				if err != nil {
					t.Errorf("TestRpcUserService_Login error %v", err)
				}
				return dal.UserDo{
					UserName: userName,
					Password: string(bytes),
				}, nil
			},
		)
		defer patches.Reset()

		// mock cache
		patches.ApplyMethod(reflect.TypeOf(s.dao),
			"SetCache",
			func(_ *dal.UserDao, _ context.Context, _ string, _ string, expire time.Duration) error {
				assert.Equal(t, expire, constant.ThreeDaysToExpire)
				return nil
			})

		patches.ApplyFunc(util.GenerateSid, func() string {
			return sessionId
		})

		//actual call
		err := s.Login(req, &resp)
		if err != nil {
			t.Errorf("TestRpcUserService_Login error %v", err)
		}

		assert.Equal(t, resp.SessionId, sessionId)
	})

	t.Run("Login no user", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(s.dao),
			"GetUserByName",
			func(_ *dal.UserDao, _ string) (dal.UserDo, error) {
				return dal.UserDo{}, errors.New("no record")
			},
		)
		defer patches.Reset()

		//actual call
		err := s.Login(req, &resp)
		if err == nil {
			t.Errorf("TestRpcUserService_Login should return error %v", err)
		}
	})

	t.Run("Login password not match", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(s.dao),
			"GetUserByName",
			func(_ *dal.UserDao, _ string) (dal.UserDo, error) {
				return dal.UserDo{
					UserName: userName,
					Password: "anyway",
				}, nil
			},
		)
		defer patches.Reset()

		//actual call
		err := s.Login(req, &resp)
		if err == nil {
			t.Errorf("TestRpcUserService_Login should return password not match %v", err)
		}

	})
}

func TestRpcUserService_GetUser(t *testing.T) {
	s := NewService(context.Background(), nil, nil)

	userName := "test"
	userId := "12345"
	nickName := "nickname"
	avatarUrl := "/test.png"
	sessionId := "testSession"
	req := entity.InfoReq{
		SessionId: sessionId,
	}
	resp := entity.UserInfoResp{}

	t.Run("GetCache from cache success", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(s), "GetSession",
			func(_ *UserService, sId string, uid *string) error {
				assert.Equal(t, sId, sessionId)
				*uid = userId
				return nil
			})
		defer patches.Reset()

		patches.ApplyMethod(reflect.TypeOf(s.dao), "GetCache",
			func(_ *dal.UserDao, _ context.Context, key string) (string, error) {
				if key == constant.UserInfoCachePrefix+userId {
					userJson, _ := json.Marshal(entity.UserInfoResp{
						UserName:  userName,
						NickName:  nickName,
						AvatarUrl: avatarUrl,
					})
					return string(userJson), nil
				} else if key == constant.SessionIdCachePrefix+sessionId {
					return userId, nil
				} else {
					return "", errors.New("not valid param")
				}

			})

		//actual call
		err := s.GetUser(req, &resp)
		if err != nil {
			t.Errorf("TestRpcUserService_GetUser error %v", err)
		}
		assert.Equal(t, resp.NickName, nickName)
		assert.Equal(t, resp.UserName, userName)
		assert.Equal(t, resp.AvatarUrl, avatarUrl)
	})

	t.Run("No session", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(s), "GetSession",
			func(_ *UserService, sId string, uid *string) error {
				return errors.New(errorcode.SessionNotExist.GetCodeStr())
			})
		defer patches.Reset()
		patches.ApplyMethod(reflect.TypeOf(s.dao), "GetCache",
			func(_ *dal.UserDao, _ context.Context, key string) (string, error) {
				if key == constant.UserInfoCachePrefix+userId {
					userJson, _ := json.Marshal(entity.UserInfoResp{
						UserName:  userName,
						NickName:  nickName,
						AvatarUrl: avatarUrl,
					})
					return string(userJson), nil
				} else if key == constant.SessionIdCachePrefix+sessionId {
					return "", errors.New(errorcode.SessionNotExist.GetCodeStr())
				} else {
					return "", errors.New("not valid param")
				}
			})

		//actual call
		err := s.GetUser(req, &resp)
		if err == nil {
			t.Errorf("TestRpcUserService_GetUser should return error %v", err)
		}
	})
}

func TestRpcUserService_EditUser(t *testing.T) {

	s := NewService(context.Background(), nil, nil)
	userId := "12345"
	userName := "test"
	nickName := "nickname"
	avatarUrl := "/test.png"
	sessionId := "testSession"

	req := entity.EditReq{
		SessionId: sessionId,
		NickName:  nickName,
		AvatarUrl: avatarUrl,
	}
	resp := entity.EditResp{}

	t.Run("Edit success", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(s.dao), "GetCache",
			func(_ *dal.UserDao, _ context.Context, key string) (string, error) {
				if key == constant.UserInfoCachePrefix+userId {
					userJson, _ := json.Marshal(entity.UserInfoResp{
						UserName:  userName,
						NickName:  nickName,
						AvatarUrl: avatarUrl,
					})
					return string(userJson), nil
				} else if key == constant.SessionIdCachePrefix+sessionId {
					fmt.Println("this is session")
					return userId, nil
				} else {
					return "", errors.New("not valid param")
				}

			})
		defer patches.Reset()

		patches.ApplyMethod(reflect.TypeOf(s.dao),
			"GetUserByUserId",
			func(_ *dal.UserDao, userId uint64) (dal.UserDo, error) {
				return dal.UserDo{}, nil
			})

		patches.ApplyMethod(reflect.TypeOf(s.dao),
			"UpdateUser",
			func(_ *dal.UserDao, _ dal.UserDo) error {
				return nil
			})

		patches.ApplyMethod(reflect.TypeOf(s.dao),
			"DeleteCache",
			func(_ *dal.UserDao, _ context.Context, key string) error {
				return nil
			})

		err := s.EditUser(req, &resp)
		if err != nil {
			t.Errorf("TestRpcUserService_EditUser error %v", err)
		}
	})

	t.Run("Edit user exists", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(s.dao), "GetCache",
			func(_ *dal.UserDao, _ context.Context, key string) (string, error) {
				if key == constant.UserInfoCachePrefix+userId {
					userJson, _ := json.Marshal(entity.UserInfoResp{
						UserName:  userName,
						NickName:  nickName,
						AvatarUrl: avatarUrl,
					})
					return string(userJson), nil
				} else if key == constant.SessionIdCachePrefix+sessionId {
					fmt.Println("this is session")
					return userId, nil
				} else {
					return "", errors.New("not valid param")
				}

			})
		defer patches.Reset()

		patches.ApplyMethod(reflect.TypeOf(s.dao),
			"GetUserByUserId",
			func(_ *dal.UserDao, userId uint64) (dal.UserDo, error) {
				return dal.UserDo{}, errors.New("user exists")
			})

		err := s.EditUser(req, &resp)
		if err == nil {
			t.Errorf("TestRpcUserService_EditUser should return user exists error %v", err)
		}
	})
}

func TestRpcUserService_UploadAvatar(t *testing.T) {
	s := NewService(context.Background(), nil, nil)

	fileName := "test.jpg"
	fileExt := ".jpg"

	fileNameMd5 := "test_md5.jpg"
	const staticUrlConst = "/static"
	const userName = "test"

	t.Run("Upload success", func(t *testing.T) {
		patches := gomonkey.ApplyFunc(util.SaveFile, func(_ *[]byte, dst string, dir string) error {
			return nil
		})
		defer patches.Reset()
		patches.ApplyFunc(util.BuildFileName, func(name string) string {
			return fileNameMd5
		})

		conf.WebConf = &conf.WebServerConf{
			StaticFileUrl: staticUrlConst,
		}
		resp := entity.UploadResp{}
		var emptyArray []byte
		err := s.UploadAvatar(entity.UploadReq{
			FileName: fileName,
			FileExt:  fileExt,
			Content:  emptyArray,
			UserName: userName,
		}, &resp)
		if err != nil {
			return
		}

		assert.Equal(t, resp.AvatarUrl, staticUrlConst+"/"+userName+"/"+fileNameMd5)
	})
}

func TestRpcUserService_GetSession(t *testing.T) {
	s := NewService(context.Background(), nil, nil)
	sessionId := "test_sessionId"
	userId := ""

	t.Run("GetCache session success", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(s.dao),
			"GetCache",
			func(_ *dal.UserDao, _ context.Context, key string) (string, error) {
				assert.Equal(t, key, constant.SessionIdCachePrefix+sessionId)
				return "12345", nil
			})
		defer patches.Reset()

		err := s.GetSession(sessionId, &userId)
		if err != nil {
			t.Errorf("GetCache session error %v", err)
			return
		}

		assert.Equal(t, userId, "12345")
	})

}
