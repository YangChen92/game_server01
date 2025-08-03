package middleware

import (
	"game/database"
	"game/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从header中获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "未提供认证令牌"})
			return
		}
		//验证Bearer格式
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "认证令牌格式错误"})
			return
		}
		tokenString := tokenParts[1]
		userID, err := utils.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "认证令牌无效"})
			return
		}

		//验证redis中token是否有效
		redisClient := database.GetRedis()
		redisToken, err := redisClient.Get("token:" + strconv.Itoa(int(userID))).Result()
		if err != nil || redisToken != tokenString {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "认证令牌无效"})
			return
		}
		c.Set("user_id", userID)
		c.Next()
		// //验证用户是否被禁用
		// user := database.GetUserByID(userID)
		// if user.Disabled {
		// 	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "用户被禁用"})
		// 	return
		// }
	}
}
