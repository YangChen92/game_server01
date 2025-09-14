package service

import (
	"encoding/json"
	"game/models"
	"game/utils"

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
		Money:        0,
		Exp:          0,
	}
	//随机生成id
	snow := &utils.Snowflake{}
	user.ID = snow.GenerateID()
	// if err := models.CreateUser(s.db, user); err != nil {
	// 	return nil, err
	// }
	userStr, err := json.Marshal(*user)
	if err != nil {
		return nil, err
	}
	//redis 所有用户信息hash存储
	if err := s.redisClient.HSet("users", fmt.Sprint(user.ID), userStr).Err(); err != nil {
		return nil, err
	}
	if err := s.redisClient.Set(utils.GetIDbyNameKey(user.Username), fmt.Sprint(user.ID), 0).Err(); err != nil {
		return nil, err
	}
	// if err := s.redisClient.HSet("name2ID", user.Username, fmt.Sprint(user.ID)).Err(); err != nil {
	// 	return nil, err
	// }
	return user, nil
}
