package category

import (
	"GoProject/model"
	"gorm.io/gorm"
)

type Repository interface {
	FindAllFinanceCategory() ([]model.FinanceCategory, error)
	CreateUserFinanceCategory(tx *gorm.DB, userFinanceCategory *model.UserFinanceCategory) error
	FindUserCategoriesByUsername(username string) ([]model.UserFinanceCategory, error)

	FindUserCategoryByUsernameAndCode(username string, code string) (model.UserFinanceCategory, error)
	CreateUserFinanceRecord(userFinanceRecord *model.UserFinanceRecord) error
	FindUserRecordsByUserIdPreload(userId uint, pageNumber int, pageSize int) (int64, []model.UserFinanceRecord, error)
	DeleteUserFinanceRecordById(recordId uint) error
	ModifyUserFinanceRecordById(recordId uint, data map[string]interface{}) error
}
