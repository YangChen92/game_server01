package service

import (
	"encoding/json"
	"game/models"
	"game/utils"
	"strconv"

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

func (s *UserService) GetUser(userId int64) (user *models.User, err error) {
	//从redis中获取用户信息
	userStr, err := s.redisClient.HGet("users", fmt.Sprint(userId)).Result()
	if err != nil {
		return nil, err
	}
	user = &models.User{}
	if err := json.Unmarshal([]byte(userStr), user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByName(username string) (user *models.User, err error) {
	//从redis中获取用户信息
	userIdStr, err := s.redisClient.Get(utils.GetIDbyNameKey(username)).Result()
	if err != nil {
		return nil, err
	}
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return nil, err
	}
	return s.GetUser(userId)
}

func (s *UserService) UpdateUser(userId int64, user *models.User) (err error) {
	//从redis中获取用户信息
	userStr, err := s.redisClient.HGet("users", fmt.Sprint(userId)).Result()
	if err != nil {
		return err
	}
	oldUser := &models.User{}
	if err := json.Unmarshal([]byte(userStr), oldUser); err != nil {
		return err
	}
	//更新用户信息
	oldUser.Username = user.Username
	oldUser.Mobile = user.Mobile
	oldUser.WechatOpenID = user.WechatOpenID
	oldUser.Password = user.Password
	//更新redis中的用户信息
	userBytes, err := json.Marshal(*oldUser)
	if err != nil {
		return err
	}
	if err := s.redisClient.HSet("users", fmt.Sprint(userId), userBytes).Err(); err != nil {
		return err
	}
	if err := s.redisClient.Set(utils.GetIDbyNameKey(oldUser.Username), fmt.Sprint(userId), 0).Err(); err != nil {
		return err
	}
	// if err := s.redisClient.HSet("name2ID", oldUser.Username, fmt.Sprint(userId)).Err(); err != nil {
	// 	return err
	// }
	//删除旧的用户名对应的id
	oldKey := utils.GetIDbyNameKey(oldUser.Username)
	if err := s.redisClient.Del(oldKey).Err(); err != nil {
		return err
	}
	return nil
}
