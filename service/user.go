package service

import (
	"encoding/json"
	"game/models"

	"fmt"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type UserService struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewUserService(db *gorm.DB, redisClient *redis.Client) *UserService {
	return &UserService{
		db:          db,
		redisClient: redisClient,
	}
}

type MobileRegisterRequest struct {
	Mobile   string
	Code     string
	Password string
}

func (s *UserService) MobileRegister(req MobileRegisterRequest) {

}

type CreateUserRequest struct {
	Username     string
	Mobile       string
	WechatOpenID string
	Password     string
}

func (s *UserService) CreateUser(req *CreateUserRequest) (user *models.User, err error) {

	user = &models.User{
		Mobile:       req.Mobile,
		Password:     req.Password,
		Username:     req.Username,
		WechatOpenID: req.WechatOpenID,
	}
	// if err := models.CreateUser(s.db, user); err != nil {
	// 	return nil, err
	// }
	userStr, err := json.Marshal(*user)
	if err != nil {
		return nil, err
	}
	//redis 所有用户信息hash存储
	if err := s.redisClient.HSet("users", "user:"+fmt.Sprint(user.ID), userStr).Err(); err != nil {
		return nil, err
	}
	return user, nil
}
