package login

import "github.com/gin-gonic/gin"

type HttpHandler interface {
	ApiLogin(c *gin.Context)
}
