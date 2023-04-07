package entity

import "gorm.io/gorm"

type TableUser struct {
	gorm.Model
	Name     string
	Password string
}

func (tableUser TableUser) TableName() string {
	return "users"
}
