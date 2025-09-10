package main

import (
	"game/controllers"
	"game/database"
	"game/middleware"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载.env文件
	_ = godotenv.Load()
	// 初始化数据库连接
	// database.InitDB()
	database.InitRedis()

	r := gin.Default()

	// 认证路由组
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/account/register", controllers.AccountRegister)
		authGroup.POST("/wechat/register", controllers.WechatRegister)
		authGroup.POST("/login", controllers.UsernameLogin)
		authGroup.POST("/wechat/login", controllers.WechatLogin)
		authGroup.POST("/sms/code", controllers.SendCaptcha)
	}

	// 需要认证的路由组
	authorized := r.Group("/")
	authorized.Use(middleware.JWTAuthMiddleware())
	{
		authorized.GET("/user/getUserInfo", controllers.GetUserInfo)
		// authorized.POST("/user/updateUserInfo", controllers.UpdateUserInfo)

	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
