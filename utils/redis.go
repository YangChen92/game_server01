package utils

import "strconv"

func GetUserKey(userId int64) string {
	return "user:" + strconv.FormatInt(userId, 10)
}

func GetCaptchaKey(mobile int64) string {
	return "captcha:" + strconv.FormatInt(mobile, 10)
}
