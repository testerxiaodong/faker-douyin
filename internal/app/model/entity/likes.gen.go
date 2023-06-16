// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package entity

import (
	"time"

	"gorm.io/gorm"
)

const TableNameLike = "likes"

// Like mapped from table <likes>
type Like struct {
	ID        int64          `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement:true" json:"id,string"`
	UserID    int64          `gorm:"column:user_id;type:int unsigned;not null" json:"user_id"`
	VideoID   int64          `gorm:"column:video_id;type:int unsigned;not null" json:"video_id"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime" json:"deleted_at"`
}

// TableName Like's table name
func (*Like) TableName() string {
	return TableNameLike
}
