package user

import "GoProject/model"

type Service interface {
	FindUser(username string) (*model.Users, error)
	FindUserByUsernamePreload(username string) (*model.Users, error)
	CreateUserData(user *model.Users, userRoleType *model.UserRoleType) error
	UpdateUser(user *model.Users, data map[string]interface{}) (*model.Users, error)
	FindUserRoleTypeByCode(code string) (*model.UserRoleType, error)
	FindUserRoleTypeByUsername(username string) (*[]model.UserRoleType, error)
	DeleteUser(user *model.Users) error
}
