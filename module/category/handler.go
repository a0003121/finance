package category

import "github.com/gin-gonic/gin"

type HttpHandler interface {
	FindUserCategories(c *gin.Context)
	CreateUserRecord(c *gin.Context)
	FindUserRecords(c *gin.Context)
}
