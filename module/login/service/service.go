package service

import (
	"GoProject/module/jwt"
	"GoProject/module/login"
	"GoProject/module/user"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	userService user.Service
}

func NewLoginService(userService user.Service) login.Service {
	return &LoginService{userService: userService}
}

func (loginService LoginService) Login(username, password string) (string, error) {
	users, userErr := loginService.userService.FindUserByUsernamePreload(username)
	if userErr != nil {
		return "", userErr
	}

	if !comparePwd(users.Password, password) {
		return "", errors.New("password error")
	}

	return jwt.GenerateToken(*users), nil
}

func comparePwd(pwd1 string, pwd2 string) bool {
	// Returns true on success, pwd1 is for the database.
	err := bcrypt.CompareHashAndPassword([]byte(pwd1), []byte(pwd2))
	if err != nil {
		return false
	} else {
		return true
	}
}
