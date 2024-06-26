package model

import "time"

func (UserFinanceCategory) TableName() string {
	return "user_finance_category"
}

type UserFinanceCategory struct {
	BaseModel
	UsersID uint
	Code    string
	Users   Users `gorm:"foreignKey:UsersID"`
}

func (FinanceCategory) TableName() string {
	return "finance_category"
}

type FinanceCategory struct {
	BaseModel
	Code string
}

func (UserFinanceRecord) TableName() string {
	return "user_finance_record"
}

type UserFinanceRecord struct {
	BaseModel
	UsersID               uint
	UserFinanceCategoryId uint
	Price                 uint
	SpendDate             time.Time
	UserFinanceCategory   UserFinanceCategory `gorm:"foreignKey:UserFinanceCategoryId"`
}
