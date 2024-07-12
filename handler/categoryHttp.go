package handler

import (
	"GoProject/common"
	"GoProject/module/category"
	"GoProject/module/user"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"log"
	"net/http"
)

type CategoryHttpHandler struct {
	categorySvc category.Service
	userSvc     user.Service
}

func NewCategoryHandler(categorySvc category.Service, userSvc user.Service, server *gin.Engine) CategoryHttpHandler {
	var handler = CategoryHttpHandler{categorySvc: categorySvc, userSvc: userSvc}

	//get user category list
	server.GET("/user/category/:username", func(c *gin.Context) {
		log.Printf("[%s]%s", c.Request.Method, c.Request.URL)
		//_ = jwt.Authenticate(c) && jwt.IsAdmin(c)
	}, func(c *gin.Context) {
		handler.findUserCategories(c)
	})

	//TODO
	//create user category
	//server.POST("/user/category/:username", func(c *gin.Context) {
	//	log.Printf("[%s]%s", c.Request.Method, c.Request.URL)
	//	//_ = jwt.Authenticate(c) && jwt.IsAdmin(c)
	//}, func(c *gin.Context) {
	//})

	//TODO
	//update user category
	//server.PUT("/user/category/:username", func(c *gin.Context) {
	//	log.Printf("[%s]%s", c.Request.Method, c.Request.URL)
	//	//_ = jwt.Authenticate(c) && jwt.IsAdmin(c)
	//}, func(c *gin.Context) {
	//})

	//TODO
	//delete user category
	//server.DELETE("/user/category/:username", func(c *gin.Context) {
	//	log.Printf("[%s]%s", c.Request.Method, c.Request.URL)
	//	//_ = jwt.Authenticate(c) && jwt.IsAdmin(c)
	//}, func(c *gin.Context) {
	//})

	return handler
}

func (handler *CategoryHttpHandler) findUserCategories(c *gin.Context) {
	var username = c.Param("username")
	result, err := handler.categorySvc.FindUserCategoriesByUsername(username)
	if err != nil {
		c.JSON(http.StatusOK, common.Fail(err.Error()))
		return
	}

	var targets []UserCategoryResponseData
	for _, userFinanceCategory := range result {
		var target UserCategoryResponseData
		copyErr := copier.Copy(&target, &userFinanceCategory)
		if copyErr != nil {
			c.JSON(http.StatusOK, common.Fail(copyErr.Error()))
			return
		}
		targets = append(targets, target)
	}
	c.JSON(http.StatusOK, common.Success(targets))
}
