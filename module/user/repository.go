package user

import (
	"GoProject/model"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(tx *gorm.DB, user *model.Users) (*model.Users, error)
	ModifyUser(user *model.Users, data map[string]interface{}) (*model.Users, error)
	FindUserByUsername(username string) (*model.Users, error)
	FindUserByUserType(userType string) (*[]model.Users, error)
	FindUserByUsernamePreload(username string) (*model.Users, error)
	DeleteUser(user *model.Users) error

	FindUserRoleTypeByCode(code string) (*model.UserRoleType, error)
	CreateUserRole(tx *gorm.DB, userRole *model.UserRole) (*model.UserRole, error)
	BeginTransaction() *gorm.DB

	FindUserRoleTypeByUsername(username string) (*[]model.UserRoleType, error)
	DeleteUserRoleByUserId(userId uint) error
}
