package delivery

import (
	"GoProject/common"
	"GoProject/module/login"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type LoginHttpHandler struct {
	loginService login.Service
}

func NewLoginHttpHandler(loginService login.Service, server *gin.Engine) login.HttpHandler {
	handler := &LoginHttpHandler{loginService: loginService}

	server.POST("/api/login", func(c *gin.Context) {
		log.Printf("[%s]%s", c.Request.Method, c.Request.URL)
		handler.ApiLogin(c)
	})

	return handler
}

func (l LoginHttpHandler) ApiLogin(c *gin.Context) {
	var username string
	var password string

	//check username and password
	if in, isExist := c.GetPostForm("username"); isExist && in != "" {
		username = in
	} else {
		c.JSON(http.StatusOK, common.Fail("必須輸入使用者名稱"))
		return
	}
	if in, isExist := c.GetPostForm("password"); isExist && in != "" {
		password = in
	} else {
		c.JSON(http.StatusOK, common.Fail("必須輸入密碼名稱"))
		return
	}

	token, err := l.loginService.Login(username, password)
	if err != nil {
		c.JSON(http.StatusOK, common.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Success(token))
}
