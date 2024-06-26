package user

import "github.com/gin-gonic/gin"

type HttpHandler interface {
	FindUser(c *gin.Context)
	CreateUsers(c *gin.Context)
	FindUserRoles(c *gin.Context)
	DeleteUser(c *gin.Context)
	UpdateUserEmail(c *gin.Context)
}
