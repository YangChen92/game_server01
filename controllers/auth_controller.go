package controllers

import (
	"encoding/json"
	"fmt"
	"game/database"
	"game/models"
	"game/service"
	"game/utils"
	redisService "game/utils/redis"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func MobileRegister(c *gin.Context) {
	type Request struct {
		Mobile   string `json:"mobile",binding:"required,len=11"`
		Code     string `json:"code",binding:"required,len=6"`
		Password string `json:"password",binding:"required,min=6,max=20"`
		Username string `json:"username",binding:"required,len=50"`
	}
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//TODO: 验证手机号码和验证码是否正确
	redisClient := database.GetRedis()
	storedCode, err := redisClient.Get("captcha:" + req.Mobile).Result()
	if err != nil || storedCode != req.Code {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误"})
		return
	}
	//TODO: 验证手机号码是否已注册

	username, err := redisClient.Get("user_iphone:" + req.Mobile).Result()
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "验证失败"})
	// 	return
	// }
	if username != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "手机号码已注册"})
		return
	}

	// db := database.GetDB()
	// if _, err := models.GetUserByMobile(db, req.Mobile); err == nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "手机号码已注册"})
	// 	return
	// }
	//密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}
	//TODO: 创建用户对象
	db := database.GetDB()
	userService := service.NewUserService(db, redisClient)
	createReq := &service.CreateUserRequest{
		Mobile:       req.Mobile,
		Password:     hashedPassword,
		Username:     req.Username,
		WechatOpenID: "",
	}
	token := ""
	if userInfo, err := userService.CreateUser(createReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	} else {
		//TODO: 注册成功返回token
		token, err = utils.GenerateToken(userInfo.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		}
	}
	// 	Mobile:   req.Mobile,
	// 	Password: hashedPassword,
	// }
	//TODO: 保存用户信息到数据库
	// db := database.GetDB()
	// if err := models.CreateUser(db, &user); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
	// 	return
	// }
	//TODO: 注册成功返回token
	c.JSON(http.StatusCreated, gin.H{"message": "注册成功", "token": token})
}

func AccountRegister(c *gin.Context) {
	type Request struct {
		Username string `json:"username",binding:"required,len=50"`
		Password string `json:"password",binding:"required,min=6,max=20"`
	}
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//TODO: 验证用户名是否已注册
	// if _, err := models.GetUserByUsername(db, req.Username); err == nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已注册"})
	// 	return
	// }
	userID, err := redisService.GetIDbyName(req.Username)
	fmt.Printf("req.userName:%s userID:%s\n", req.Username, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或密码错误"})
		return
	}
	if userID != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已注册"})
		return
	}

	//密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}
	//TODO: 创建用户对象
	redisClient := database.GetRedis()
	db := database.GetDB()
	userService := service.NewUserService(db, redisClient)
	createReq := &service.CreateUserRequest{
		Mobile:       "",
		Password:     hashedPassword,
		Username:     req.Username,
		WechatOpenID: "",
	}
	token := ""
	if userInfo, err := userService.CreateUser(createReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	} else {
		//TODO: 注册成功返回token
		token, err = utils.GenerateToken(userInfo.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		}
		//保存token到redis
		redisClient.Set("token:"+token, userInfo.ID, 24*time.Hour)
	}
	// //TODO: 创建用户对象
	// user := models.User{
	// 	Username: req.Username,
	// 	Password: hashedPassword,
	// }
	// //TODO: 保存用户信息到数据库
	// if err := models.CreateUser(db, &user); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
	// 	return
	// }
	//TODO: 注册成功返回token
	c.JSON(http.StatusCreated, gin.H{"message": "注册成功", "token": token})
}

// 微信小程序注册
func WechatRegister(c *gin.Context) {
	type Request struct {
		WechatOpenID string `json:"wechat_openid",binding:"required,len=100"`
		Username     string `json:"username",binding:"required,len=50"`
		Password     string `json:"password",binding:"required,min=6,max=20"`
		Mobile       string `json:"mobile",binding:"required,len=11"`
	}
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//TODO: 验证微信openid是否已注册
	db := database.GetDB()
	if _, err := models.GetUserByWechatOpenID(db, req.WechatOpenID); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "微信openid已注册"})
		return
	}

	//检查用户名
	result := db.Table("users").Where("username = ?", req.Username).First(&models.User{})
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	//密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}
	//TODO: 创建用户对象
	// db := database.GetDB()
	redisClient := database.GetRedis()
	userService := service.NewUserService(db, redisClient)
	createReq := &service.CreateUserRequest{
		Mobile:       req.Mobile,
		Password:     hashedPassword,
		Username:     req.Username,
		WechatOpenID: req.WechatOpenID,
	}
	token := ""
	strUserInfo := []byte("")
	if userInfo, err := userService.CreateUser(createReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	} else {
		token, err = utils.GenerateToken(userInfo.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
			return
		}
		strUserInfo, err = json.Marshal(userInfo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "strUserInfo error"})
			return
		}
	}
	// //TODO: 创建用户对象
	// user := models.User{
	// 	WechatOpenID: req.WechatOpenID,
	// 	Username:     req.Username,
	// 	Password:     hashedPassword,
	// 	Mobile:       req.Mobile,
	// }
	// //TODO: 保存用户信息到数据库
	// if err := models.CreateUser(db, &user); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
	// 	return
	// }
	//TODO: 注册成功返回token

	c.JSON(http.StatusCreated, gin.H{"message": "注册成功", "token": token, "userInfo": strUserInfo})

}

type LoginRS struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	Exp          uint64 `json:"exp"`
	Money        uint64 `json:"money"`
	WechatOpenID string `json:"wechat_openid"`
}

// 用户名登录
func UsernameLogin(c *gin.Context) {
	type Request struct {
		Username string `json:"username",binding:"required,len=50"`
		Password string `json:"password",binding:"required,min=6,max=20"`
	}
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//TODO: 验证用户名和密码是否正确
	db := database.GetDB()
	user, err := models.GetUserByUsername(db, req.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或密码错误"})
		return
	}
	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或密码错误"})
		return
	}
	//TODO: 登录成功返回token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}
	//保存token到redis
	redisClient := database.GetRedis()
	redisClient.Set("token:"+token, user.ID, 24*time.Hour)
	c.JSON(http.StatusOK, gin.H{"message": "登录成功", "user_id": user.ID, "token": token})
}

func WechatLogin(c *gin.Context) {
	type Request struct {
		OpenID string `json:"openid",binding:"required,len=100"`
	}
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//TODO: 验证微信openid是否已注册
	db := database.GetDB()
	user, err := models.GetUserByWechatOpenID(db, req.OpenID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "微信openid未注册"})
		return
	}
	//TODO: 登录成功返回token
	c.JSON(http.StatusOK, gin.H{"message": "登录成功", "user_id": user.ID})
}

// 发送验证码
func SendCaptcha(c *gin.Context) {
	type Request struct {
		Mobile string `json:"mobile",binding:"required,len=11"`
	}
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//TODO: 验证手机号码是否已注册
	db := database.GetDB()
	if _, err := models.GetUserByMobile(db, req.Mobile); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "手机号码已注册"})
		return
	}
	//TODO: 生成验证码
	code := "123456"
	//TODO: 保存验证码到redis
	redisClient := database.GetRedis()
	redisClient.Set("captcha:"+req.Mobile, code, 10*time.Minute)
	//TODO: 发送验证码
	//TODO: 返回成功信息
	c.JSON(http.StatusOK, gin.H{"message": "验证码发送成功"})
}
