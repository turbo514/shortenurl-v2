package model

import (
	"gorm.io/gorm"
	"time"
)

type Tenant struct {
	ID        []byte         `gorm:"primaryKey;column:id;type:binary(16)"`
	Name      string         `gorm:"column:name"`
	ApiKey    string         `gorm:"column:api_key"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (tenant *Tenant) TableName() string {
	return "tenants"
}
