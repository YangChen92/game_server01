package redisService

import (
	"encoding/json"
	"game/database"
	"game/models"
	"game/utils"

	"github.com/go-redis/redis"
)

func GetRedisAuth(userId int64) (userData *models.User, err error) {
	key := utils.GetUserKey(userId)
	cli := database.GetRedis()
	UserValue, err := cli.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	userData = &models.User{}
	err = json.Unmarshal([]byte(UserValue), &userData)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

func SetRedisAuth(userData *models.User) error {
	key := utils.GetUserKey(userData.ID)
	cli := database.GetRedis()
	UserValue, err := json.Marshal(userData)
	if err != nil {
		return err
	}
	err = cli.Set(key, UserValue, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func DelRedisAuth(userId int64) error {
	key := utils.GetUserKey(userId)
	cli := database.GetRedis()
	err := cli.Del(key).Err()
	if err != nil {
		return err
	}
	return nil
}

func UpdateRedisAuth(userData *models.User) error {
	key := utils.GetUserKey(userData.ID)
	cli := database.GetRedis()
	UserValue, err := json.Marshal(userData)
	if err != nil {
		return err
	}
	err = cli.Set(key, UserValue, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func SetIDbyName(name string) (int64, error) {
	key := utils.GetIDbyNameKey(name)
	cli := database.GetRedis()
	id, err := cli.Incr(key).Result()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetIDbyName(name string) (string, error) {
	key := utils.GetIDbyNameKey(name)
	cli := database.GetRedis()
	id, err := cli.Get(key).Result()

	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return id, nil
}

// ==============================
