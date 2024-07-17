package category

import (
	"GoProject/model"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	BeginTransaction() *gorm.DB
	FindAllFinanceCategory() ([]model.FinanceCategory, error)
	CreateUserFinanceCategory(tx *gorm.DB, userFinanceCategory *model.UserFinanceCategory) error
	FindUserCategoriesByUsername(username string) ([]model.UserFinanceCategory, error)

	FindUserCategoryByUsernameAndCode(username string, code string) (model.UserFinanceCategory, error)

	CreateUserFinanceRecord(userFinanceRecord *model.UserFinanceRecord) error
	CreateUserFinanceRecords(userFinanceRecord *[]model.UserFinanceRecord) error
	FindUserRecordsByUserIdPreload(userId uint, pageNumber int, pageSize int, startDateTime time.Time, endDateTime time.Time) (int64, []model.UserFinanceRecord, error)
	DeleteUserFinanceRecordById(recordId uint) error
	ModifyUserFinanceRecordById(recordId uint, data map[string]interface{}) error
	FindUserRecordsByUsernamePreload(username string) ([]model.UserFinanceRecord, error)
	ModifyUserFinanceCategory(userFinanceCategory *model.UserFinanceCategory, data map[string]interface{}) error
	DeleteUserFinanceRecordByCategoryId(tx *gorm.DB, categoryId uint) error
	DeleteUserFinanceCategoryById(tx *gorm.DB, categoryId uint) error
}
