package service

import (
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
