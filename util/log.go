// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/6/13

package util

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

var Logger = log.New()

func InitLog(logPath string, level log.Level) {
	Logger.SetFormatter(&log.JSONFormatter{})

	Logger.SetLevel(level)

	if strings.ToLower(logPath) == "stdout" {
		Logger.SetOutput(os.Stdout)
	} else {
		Logger.SetOutput(&lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     7,    //days
			Compress:   true, // disabled by default
		})
	}
}
