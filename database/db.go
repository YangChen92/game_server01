package database

import (
	"game/config"

	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var rdb *redis.Client

func InitDB() {
	dsn := config.GetMysqlDSN()
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database, error: " + err.Error())
	}
}

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.GetRedisAddr(),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if err := rdb.Ping().Err(); err != nil {
		panic("failed to connect redis, error: " + err.Error())
	}
}

func GetDB() *gorm.DB {
	return db
}

func GetRedis() *redis.Client {
	return rdb
}
