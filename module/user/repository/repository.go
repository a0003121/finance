package repository

import (
	. "GoProject/model"
	"GoProject/module/user"
	"gorm.io/gorm"
)

type UserRepository struct {
	orm *gorm.DB
}

func NewUserRepository(orm *gorm.DB) user.Repository {
	return &UserRepository{orm: orm}
}

func (u UserRepository) FindUserByUsernamePreload(username string) (*Users, error) {
	var users Users
	result := u.orm.Preload("UserRoles.UserRoleType").First(&users, "username = ?", username)
	return &users, result.Error
}

func (u UserRepository) DeleteUser(user *Users) error {
	err := u.orm.Delete(&user).Error
	return err
}

func (u UserRepository) FindUserRoleTypeByUsername(username string) (*[]UserRoleType, error) {
	var userRoleTypes []UserRoleType
	err := u.orm.Joins("JOIN user_roles ON user_role_type.id = user_roles.user_role_type_id").
		Joins("JOIN users ON users.id = user_roles.users_id").
		Where("users.username = ?", username).
		Find(&userRoleTypes).Error
	return &userRoleTypes, err
}

func (u UserRepository) DeleteUserRoleByUserId(userId uint) error {
	return u.orm.Where("users_id = ?", userId).Delete(&UserRole{}).Error

}

func (u UserRepository) BeginTransaction() *gorm.DB {
	return u.orm.Begin()
}

func (u UserRepository) CreateUser(tx *gorm.DB, user *Users) (*Users, error) {
	var err error
	if tx != nil {
		err = tx.Create(&user).Error
	} else {
		err = u.orm.Create(&user).Error
	}
	return user, err
}

func (u UserRepository) ModifyUser(user *Users, data map[string]interface{}) (*Users, error) {
	var err error
	err = u.orm.Model(&user).Updates(data).Error
	return user, err
}

func (u UserRepository) FindUserByUsername(username string) (*Users, error) {
	var users *Users
	result := u.orm.Where("username = ?", username).First(&users)
	return users, result.Error
}

func (u UserRepository) CreateUserRole(tx *gorm.DB, userRole *UserRole) (*UserRole, error) {
	var err error

	if tx != nil {
		err = tx.Create(&userRole).Error
	} else {
		err = u.orm.Create(&userRole).Error
	}
	return userRole, err
}

func (u UserRepository) FindUserRoleTypeByCode(code string) (*UserRoleType, error) {
	var userRoleType *UserRoleType
	roleTypeResult := u.orm.Where("code = ?", code).First(&userRoleType)
	return userRoleType, roleTypeResult.Error
}
