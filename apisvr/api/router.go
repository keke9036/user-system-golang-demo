package api

import (
	"entry-task/bizsvr/constant"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	// bean url-handler mapping
	user := NewUser()
	loginApiGroup := r.Group("/api/v1/user")
	{
		loginApiGroup.POST("/login", user.Login)
		loginApiGroup.POST("/register", user.Register)
		loginApiGroup.GET("/ping", user.Ping)
	}

	// need login
	sessionAipGroup := r.Group("/api/v1/user")
	sessionAipGroup.Use(SessionRequired)
	{
		sessionAipGroup.POST("/edit", user.Edit)
		sessionAipGroup.POST("/uploadAvatar", user.UploadAvatar)
		sessionAipGroup.POST("/logout", user.Logout)

		sessionAipGroup.GET("/info", user.Info)
	}

	// static file
	r.Static("/upload", constant.UploadFileDir)

	r.Static("/static", "./static")
	r.StaticFile("/index", "./static/html/login.html")

	return r
}
