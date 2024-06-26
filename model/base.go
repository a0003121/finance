package model

import "time"

type BaseModel struct {
	ID         uint      `json:"id" gorm:"primary_key;auto_increase'"`
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime time.Time `json:"updateTime" gorm:"autoUpdateTime"`
}
