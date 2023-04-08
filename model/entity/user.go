package entity

import "gorm.io/gorm"

type TableUser struct {
	gorm.Model
	Name     string
	Password string
}

// TableName 修改映射的表名为users
func (tableUser TableUser) TableName() string {
	return "users"
}
