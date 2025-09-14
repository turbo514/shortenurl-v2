package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        []byte         `gorm:"primaryKey;column:id"`
	TenantID  []byte         `gorm:"column:tenant_id"`
	Name      string         `gorm:"column:name"`
	Password  string         `gorm:"column:password"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (u *User) TableName() string {
	return "users"
}
