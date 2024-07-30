package service

import (
	"GoProject/model"
	"GoProject/module/category"
	"gorm.io/gorm"
	"time"
)

type CategoryService struct {
	repo category.Repository
}

func NewCategoryService(categoryRepository category.Repository) category.Service {
	return &CategoryService{repo: categoryRepository}
}

func (c CategoryService) CreateAllUserFinanceCategory(tx *gorm.DB, users model.Users) error {
	defaultCategories, findError := c.repo.FindAllFinanceCategory()
	if findError != nil {
		return findError
	}
	for _, financeCategory := range defaultCategories {
		userCategory := model.UserFinanceCategory{UsersID: users.ID, Code: financeCategory.Code}
		createErr := c.repo.CreateUserFinanceCategory(tx, &userCategory)
		if createErr != nil {
			return createErr
		}
	}
	return nil
}

func (c CategoryService) DeleteUserFinanceCategoryAndRelatedRecord(userFinanceCategory model.UserFinanceCategory) error {
	tx := c.repo.BeginTransaction()
	deleteError2 := c.repo.DeleteUserFinanceRecordByCategoryId(tx, userFinanceCategory.ID)
	if deleteError2 != nil {
		return deleteError2
	}
	deleteError1 := c.repo.DeleteUserFinanceCategoryById(tx, userFinanceCategory.ID)
	if deleteError1 != nil {
		return deleteError1
	}
	return tx.Commit().Error
}

func (c CategoryService) FindUserCategoriesByUsername(username string) ([]model.UserFinanceCategory, error) {
	return c.repo.FindUserCategoriesByUsername(username)
}

func (c CategoryService) FindUserCategoryByUsernameAndCode(username string, code string) (model.UserFinanceCategory, error) {
	return c.repo.FindUserCategoryByUsernameAndCode(username, code)
}

func (c CategoryService) CreateUserFinanceRecord(userFinanceRecord *model.UserFinanceRecord) error {
	return c.repo.CreateUserFinanceRecord(userFinanceRecord)
}

func (c CategoryService) CreateUserFinanceRecords(userFinanceRecord *[]model.UserFinanceRecord) error {
	return c.repo.CreateUserFinanceRecords(userFinanceRecord)
}

func (c CategoryService) FindUserRecordsPageByUserIdPreload(userId uint, pageNumber int, pageSize int, startDateTime time.Time, endDateTime time.Time) (int64, []model.UserFinanceRecord, error) {
	return c.repo.FindUserRecordsPageByUserIdPreload(userId, pageNumber, pageSize, startDateTime, endDateTime)
}

func (c CategoryService) FindUserRecordsByUserIdPreload(userId uint, startDateTime time.Time, endDateTime time.Time) ([]model.UserFinanceRecord, error) {
	return c.repo.FindUserRecordsByUserIdPreload(userId, startDateTime, endDateTime)
}

func (c CategoryService) DeleteUserFinanceRecordById(recordId uint) error {
	return c.repo.DeleteUserFinanceRecordById(recordId)
}

func (c CategoryService) ModifyUserFinanceRecordById(recordId uint, data map[string]interface{}) error {
	return c.repo.ModifyUserFinanceRecordById(recordId, data)
}

func (c CategoryService) ModifyUserFinanceCategory(userFinanceCategory *model.UserFinanceCategory, data map[string]interface{}) error {
	return c.repo.ModifyUserFinanceCategory(userFinanceCategory, data)
}

func (c CategoryService) FindUserRecordsByUsernamePreload(username string) ([]model.UserFinanceRecord, error) {
	return c.repo.FindUserRecordsByUsernamePreload(username)
}

func (c CategoryService) CreateUserFinanceCategory(financeCategory *model.UserFinanceCategory) error {
	return c.repo.CreateUserFinanceCategory(nil, financeCategory)
}
