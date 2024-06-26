package category

import (
	"GoProject/model"
	"gorm.io/gorm"
)

type Service interface {
	CreateUserFinanceCategory(tx *gorm.DB, users model.Users) error
	FindUserCategoriesByUsername(username string) ([]model.UserFinanceCategory, error)

	FindUserCategoryByUsernameAndCode(username string, code string) (model.UserFinanceCategory, error)
	CreateUserFinanceRecord(userFinanceRecord *model.UserFinanceRecord) error
	FindUserRecordsByUserIdPreload(userId uint, pageNumber int, pageSize int) (int64, []model.UserFinanceRecord, error)
}