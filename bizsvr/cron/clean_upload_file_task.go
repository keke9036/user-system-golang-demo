// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/7/3

package cron

import (
	"entry-task/bizsvr/bean"
	"entry-task/bizsvr/constant"
	"entry-task/bizsvr/dal"
	"entry-task/util"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func RunCleanTask() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for {
			<-ticker.C
			doClean()
		}
	}()
}

func doClean() {
	files, err := ioutil.ReadDir(constant.UploadFileDir)
	if err != nil {
		util.Logger.Errorf("ReadDir error: %s", constant.UploadFileDir)
		return
	}

	for _, user := range files {
		if !user.IsDir() {
			util.Logger.Errorf("Not dir %s", user.Name())
			continue
		}

		images, err := ioutil.ReadDir(constant.UploadFileDir + "/" + user.Name())
		if err != nil {
			util.Logger.Errorf("ReadDir error: %s", user.Name())
			return
		}

		dao := dal.NewUserDao(bean.Db, nil)
		userDo, err := dao.GetUserByName(user.Name())
		if err != nil {
			util.Logger.Errorf("Query user %s err, %v", user.Name(), err)
			continue
		}
		if !userDo.AvatarUrl.Valid {
			continue
		}

		baseDir := constant.UploadFileDir + "/" + user.Name()
		url := userDo.AvatarUrl.String
		index := strings.LastIndex(url, "/")
		fileName := url[index+1:]
		for _, image := range images {
			if image.Name() != fileName {
				util.Logger.Infof("Delete file %s", image.Name())
				err := os.Remove(baseDir + "/" + image.Name())
				if err != nil {
					util.Logger.Errorf("Remove file %s error %v", image.Name(), err)
					return
				}
			}
		}
	}

}
