// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/6/17

package dal

import (
	"entry-task/conf"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/magiconair/properties/assert"
	"testing"
)

var userDao *UserDao

func setup(t *testing.T) {
	err := conf.LoadConf()
	if err != nil {
		t.Errorf("Loadconf error, %v", err)
		return
	}
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		conf.DbConf.Username,
		conf.DbConf.Password,
		conf.DbConf.Host,
		conf.DbConf.Port,
		conf.DbConf.DbName)
	db, err := sqlx.Connect("mysql", dbUrl)
	if err != nil {
		t.Errorf("mysql connect error %v", err)
		return
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	userDao = NewUserDao(db, nil)
}

func TestUserDao_QueryUserByName(t *testing.T) {
	setup(t)
	userName := "testa"

	userDo, err := userDao.GetUserByName(userName)
	if err != nil {
		t.Errorf("TestUserDao_QueryUserByName error, %v", err)
		return
	}

	assert.Equal(t, userDo.NickName, "测试账号")
}

func TestUserDao_QueryUserByUserId(t *testing.T) {
	setup(t)
	userId := 283578804197265408

	userDo, err := userDao.GetUserByUserId(uint64(userId))
	if err != nil {
		t.Errorf("TestUserDao_QueryUserByUserId error, %v", err)
		return
	}

	assert.Equal(t, userDo.NickName, "拉克斯")
}

func TestUserDao_UpdateUser(t *testing.T) {
	setup(t)
	userId := 283578804197265408
	nickName := "拉克斯.测试"

	toUpdateUser := UserDo{
		UserId:   uint64(userId),
		NickName: nickName,
	}
	err := userDao.UpdateUser(toUpdateUser)
	if err != nil {
		t.Errorf("TestUserDao_UpdateUser error, %v", err)
		return
	}
	userDo, err := userDao.GetUserByUserId(uint64(userId))
	if err != nil {
		t.Errorf("TestUserDao_QueryUserByName error, %v", err)
		return
	}

	assert.Equal(t, userDo.NickName, nickName)

}
