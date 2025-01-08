package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db}

}
func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.UpdateTime = now
	u.CreateTime = now
	return dao.db.WithContext(ctx).Create(&u).Error
}

// User 直接对应数据库表结构
// 有些人叫做 entity，有些人叫做model，有些人叫做 PO(persistent object)
type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	// 创建时间，毫秒数
	CreateTime int64
	// 更新时间，毫秒数
	UpdateTime int64
}
