package service

import (
	. "GoProject/model"
	"GoProject/module/category"
	"GoProject/module/user"
)

type UserService struct {
	repo            user.Repository
	categoryService category.Service
}

func NewUserService(repo user.Repository, service category.Service) user.Service {
	return &UserService{repo: repo, categoryService: service}
}

func (u UserService) FindUserByUsernamePreload(username string) (*Users, error) {
	return u.repo.FindUserByUsernamePreload(username)
}

func (u UserService) UpdateUser(user *Users, data map[string]interface{}) (*Users, error) {
	return u.repo.ModifyUser(user, data)
}

func (u UserService) DeleteUser(user *Users) error {
	tx := u.repo.BeginTransaction()

	userRoleErr := u.repo.DeleteUserRoleByUserId(user.ID)
	if userRoleErr != nil {
		tx.Rollback()
		return userRoleErr
	}

	userErr := u.repo.DeleteUser(user)
	if userErr != nil {
		tx.Rollback()
		return userErr
	}

	return tx.Commit().Error
}

func (u UserService) FindUserRoleTypeByUsername(username string) (*[]UserRoleType, error) {
	return u.repo.FindUserRoleTypeByUsername(username)
}

func (u UserService) FindUserRoleTypeByCode(code string) (*UserRoleType, error) {
	return u.repo.FindUserRoleTypeByCode(code)
}

func (u UserService) CreateUserData(user *Users, userRoleType *UserRoleType) error {
	tx := u.repo.BeginTransaction()

	_, userErr := u.repo.CreateUser(tx, user)
	if userErr != nil {
		tx.Rollback()
		return userErr
	}

	userRole := UserRole{
		UserRoleTypeID: userRoleType.ID,
		UsersID:        user.ID,
	}
	_, userRoleErr := u.repo.CreateUserRole(tx, &userRole)
	if userRoleErr != nil {
		tx.Rollback()
		return userErr
	}

	categoryErr := u.categoryService.CreateUserFinanceCategory(tx, *user)
	if categoryErr != nil {
		tx.Rollback()
		return categoryErr
	}

	return tx.Commit().Error
}

func (u UserService) FindUser(username string) (*Users, error) {
	return u.repo.FindUserByUsername(username)
}
