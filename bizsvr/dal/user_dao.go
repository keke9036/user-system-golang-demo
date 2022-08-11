// @Description user related dal, including RDS&Redis
// @Author weitao.yin@shopee.com
// @Since 2022/6/13

package dal

import (
	"context"
	"database/sql"
	errorcode "entry-task/error"
	"entry-task/util"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"time"
)

type UserDao struct {
	db    *sqlx.DB
	cache *redis.Client
}

func NewUserDao(db *sqlx.DB, cache *redis.Client) *UserDao {
	return &UserDao{db: db, cache: cache}
}

type UserDo struct {
	UserId     uint64         `db:"user_id"`
	UserName   string         `db:"user_name"`
	Password   string         `db:"password"`
	NickName   string         `db:"nick_name"`
	AvatarUrl  sql.NullString `db:"avatar_url"`
	CreateTime int64          `db:"create_time"`
	ModifyTime int64          `db:"modify_time"`
}

func (d *UserDao) GetUserByName(userName string) (UserDo, error) {
	if len(userName) <= 0 {
		return UserDo{}, errors.New(errorcode.ParamInvalid.Msg)
	}
	sql := "select user_id, user_name, password, nick_name, avatar_url, create_time, modify_time from user_tab where user_name = ?"
	var u UserDo
	err := d.db.Get(&u, sql, userName)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (d *UserDao) GetUserByUserId(userId uint64) (UserDo, error) {
	sql := "select user_id, user_name, password, nick_name, avatar_url, create_time, modify_time from user_tab where user_id = ?"
	var u UserDo
	err := d.db.Get(&u, sql, userId)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (d *UserDao) UpdateUser(u UserDo) error {
	if (UserDo{}) == u {
		return errors.New(errorcode.ParamInvalid.Msg)
	}
	now := time.Now().UnixMilli()
	u.ModifyTime = now
	sql := "update user_tab set nick_name = ifnull(nullif(?, ''), nick_name), avatar_url = ifnull(nullif(?, ''), avatar_url), modify_time = ? where user_id = ?"
	res, err := d.db.Exec(sql, u.NickName, u.AvatarUrl.String, u.ModifyTime, u.UserId)
	if err != nil {
		return err
	}

	row, err := res.RowsAffected()
	if err != nil {
		return err
	}
	util.Logger.Infof("Update user %d, nickName %s, avatarUrl %s, affect %d row",
		u.UserId, u.NickName, u.AvatarUrl.String, row)

	return nil
}

func (d *UserDao) AddUser(u UserDo) error {
	if (UserDo{}) == u {
		return errors.New(errorcode.ParamInvalid.Msg)
	}
	now := time.Now().UnixMilli()
	u.CreateTime = now
	u.ModifyTime = now
	sql := "insert into user_tab(user_id, user_name, password, nick_name, create_time, modify_time) values (?, ?, ?, ?, ?, ?)"

	res, err := d.db.Exec(sql, u.UserId, u.UserName, u.Password, u.NickName, u.CreateTime, u.ModifyTime)
	if err != nil {
		return err
	}

	row, err := res.RowsAffected()
	if err != nil {
		return err
	}
	util.Logger.Infof("Create user %s, nickName %s, affect %d row",
		u.UserName, u.NickName, row)
	return nil
}

func (d *UserDao) GetCache(ctx context.Context, key string) (string, error) {
	value, err := d.cache.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (d *UserDao) SetCache(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := d.cache.Set(ctx, key, value, expiration).Err()
	return err
}

func (d *UserDao) DeleteCache(ctx context.Context, key string) error {
	err := d.cache.Del(ctx, key).Err()
	return err
}
