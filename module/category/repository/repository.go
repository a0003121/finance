package repository

import (
	"GoProject/model"
	"GoProject/module/category"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	orm *gorm.DB
}

func NewCategoryRepository(orm *gorm.DB) category.Repository {
	return &CategoryRepository{orm: orm}
}

func (c CategoryRepository) FindAllFinanceCategory() ([]model.FinanceCategory, error) {
	var result []model.FinanceCategory
	err := c.orm.Find(&result)
	return result, err.Error
}

func (c CategoryRepository) CreateUserFinanceCategory(tx *gorm.DB, userFinanceCategory *model.UserFinanceCategory) error {
	if tx != nil {
		return tx.Create(&userFinanceCategory).Error
	}
	return c.orm.Create(&userFinanceCategory).Error
}

func (c CategoryRepository) FindUserCategoriesByUsername(username string) ([]model.UserFinanceCategory, error) {
	var result []model.UserFinanceCategory

	err := c.orm.Table("user_finance_category").
		Select("user_finance_category.*").
		Joins("JOIN users ON users.id = user_finance_category.users_id").
		Where("users.username = ?", username).
		Find(&result).Error

	return result, err
}

func (c CategoryRepository) FindUserCategoryByUsernameAndCode(username string, code string) (model.UserFinanceCategory, error) {
	var result model.UserFinanceCategory
	err := c.orm.Table("user_finance_category").
		Select("user_finance_category.*").
		Joins("JOIN users ON users.id = user_finance_category.users_id").
		Where("users.username = ? and user_finance_category.code = ?", username, code).
		First(&result).Error

	return result, err
}

func (c CategoryRepository) CreateUserFinanceRecord(userFinanceRecord *model.UserFinanceRecord) error {
	return c.orm.Create(&userFinanceRecord).Error
}

func (c CategoryRepository) FindUserRecordsByUserIdPreload(userId uint, pageNumber int, pageSize int) (int64, []model.UserFinanceRecord, error) {
	var result []model.UserFinanceRecord
	offset := (pageNumber - 1) * pageSize

	var count int64
	countErr := c.orm.Table("user_finance_record").
		Where("users_id = ?", userId).
		Count(&count).
		Error

	if countErr != nil {
		return 0, nil, countErr
	}

	err := c.orm.Preload("UserFinanceCategory").
		Where("users_id = ?", userId).
		Order("spend_date desc").
		Offset(offset). // Offset for the pages
		Limit(pageSize). // Limit for the page size
		Find(&result).
		Error
	return count, result, err
}

func (c CategoryRepository) DeleteUserFinanceRecordById(recordId uint) error {
	var userFinanceRecord model.UserFinanceRecord
	return c.orm.Delete(&userFinanceRecord, recordId).Error
}

func (c CategoryRepository) ModifyUserFinanceRecordById(recordId uint, data map[string]interface{}) error {
	var err error
	var userFinanceRecord model.UserFinanceRecord
	findErr := c.orm.First(&userFinanceRecord, recordId).Error

	if findErr != nil {
		return findErr
	}
	err = c.orm.Model(&userFinanceRecord).Updates(data).Error
	return err
}
