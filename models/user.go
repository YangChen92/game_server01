package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint   `gorm:"primary_key"`
	Username     string `gorm:"size:50;unique_index"`
	Mobile       string `gorm:"size:20;unique_index"`
	WechatOpenID string `gorm:"size:100;unique_index"`
	Password     string `gorm:"size:100"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *User) TableName() string {
	return "user"
}

func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	var user User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func GetUserByMobile(db *gorm.DB, mobile string) (*User, error) {
	var user User
	if err := db.Where("mobile = ?", mobile).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func GetUserByWechatOpenID(db *gorm.DB, wechatOpenID string) (*User, error) {
	var user User
	if err := db.Where("wechat_open_id = ?", wechatOpenID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
