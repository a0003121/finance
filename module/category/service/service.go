package service

import (
	"GoProject/model"
	"GoProject/module/category"
	"gorm.io/gorm"
)

type CategoryService struct {
	repo category.Repository
}

func NewCategoryService(categoryRepository category.Repository) category.Service {
	return &CategoryService{repo: categoryRepository}
}

func (c CategoryService) CreateUserFinanceCategory(tx *gorm.DB, users model.Users) error {
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

func (c CategoryService) FindUserCategoriesByUsername(username string) ([]model.UserFinanceCategory, error) {
	return c.repo.FindUserCategoriesByUsername(username)
}

func (c CategoryService) FindUserCategoryByUsernameAndCode(username string, code string) (model.UserFinanceCategory, error) {
	return c.repo.FindUserCategoryByUsernameAndCode(username, code)
}

func (c CategoryService) CreateUserFinanceRecord(userFinanceRecord *model.UserFinanceRecord) error {
	return c.repo.CreateUserFinanceRecord(userFinanceRecord)
}

func (c CategoryService) FindUserRecordsByUserIdPreload(userId uint, pageNumber int, pageSize int) (int64, []model.UserFinanceRecord, error) {
	return c.repo.FindUserRecordsByUserIdPreload(userId, pageNumber, pageSize)
}
