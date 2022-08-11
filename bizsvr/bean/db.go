// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/6/13

package bean

import (
	"entry-task/conf"
	"entry-task/util"
	"fmt"
	"github.com/jmoiron/sqlx"
)

var (
	Db *sqlx.DB
)

func InitDb(conf *conf.DbServerConf) error {
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DbName)
	db, err := sqlx.Connect("mysql", dbUrl)
	if err != nil {
		util.Logger.Fatalf("Init db fail, dbUrl: %s", dbUrl)
		return err
	}

	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(100)
	Db = db

	return nil
}
