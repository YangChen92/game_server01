package controllers

import "github.com/gin-gonic/gin"

// func Handle
func GetUserInfo(ctx *gin.Context) {
	// TODO: get user info by userId
	userId := ctx.Param("user_id")
	if userId == "" {
		ctx.JSON(400, gin.H{
			"error": "user_id is required",
		})
		return
	}
	//从redis中获取用户信息
}
