package handler

import (
	"GoProject/common"
	"GoProject/model"
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

	//create user category
	server.POST("/user/category/:username", func(c *gin.Context) {
		log.Printf("[%s]%s", c.Request.Method, c.Request.URL)
		//_ = jwt.Authenticate(c) && jwt.IsAdmin(c)
	}, func(c *gin.Context) {
		handler.createUserFinanceCategory(c)
	})

	//update user category
	server.PUT("/user/category/:username", func(c *gin.Context) {
		log.Printf("[%s]%s", c.Request.Method, c.Request.URL)
		//_ = jwt.Authenticate(c) && jwt.IsAdmin(c)
	}, func(c *gin.Context) {
		handler.updateUserFinanceCategory(c)
	})

	//delete user category
	server.DELETE("/user/category/:username", func(c *gin.Context) {
		log.Printf("[%s]%s", c.Request.Method, c.Request.URL)
		//_ = jwt.Authenticate(c) && jwt.IsAdmin(c)
	}, func(c *gin.Context) {
		handler.deleteUserFinanceCategory(c)
	})

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

func (handler *CategoryHttpHandler) createUserFinanceCategory(c *gin.Context) {
	var username = c.Param("username")
	var requestBody CreateUserFinanceCategoryRequestData

	if err := c.ShouldBindBodyWithJSON(&requestBody); err != nil {
		c.JSON(http.StatusOK, common.Fail(err.Error()))
		return
	}

	users, userErr := handler.userSvc.FindUser(username)
	if userErr != nil {
		c.JSON(http.StatusOK, common.Fail(userErr.Error()))
		return
	}
	createObj := model.UserFinanceCategory{
		UsersID: users.ID,
		Code:    requestBody.Code,
	}
	findErr := handler.categorySvc.CreateUserFinanceCategory(&createObj)
	if findErr != nil {
		c.JSON(http.StatusOK, common.Fail(findErr.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Success(""))
}

type CreateUserFinanceCategoryRequestData struct {
	Code string `json:"code" binding:"required"`
}

func (handler *CategoryHttpHandler) updateUserFinanceCategory(c *gin.Context) {
	var username = c.Param("username")
	var requestBody UpdateUserFinanceCategoryRequestData

	if err := c.ShouldBindBodyWithJSON(&requestBody); err != nil {
		c.JSON(http.StatusOK, common.Fail(err.Error()))
		return
	}

	users, userErr := handler.categorySvc.FindUserCategoryByUsernameAndCode(username, requestBody.OldCode)
	if userErr != nil {
		c.JSON(http.StatusOK, common.Fail(userErr.Error()))
		return
	}

	datas := map[string]interface{}{
		"Code": requestBody.NewCode,
	}
	updateErr := handler.categorySvc.ModifyUserFinanceCategory(&users, datas)
	if updateErr != nil {
		c.JSON(http.StatusOK, common.Fail(updateErr.Error()))
	}

	c.JSON(http.StatusOK, common.Success(""))
}

type UpdateUserFinanceCategoryRequestData struct {
	OldCode string `json:"old_code" binding:"required"`
	NewCode string `json:"new_code" binding:"required"`
}

func (handler *CategoryHttpHandler) deleteUserFinanceCategory(c *gin.Context) {
	var username = c.Param("username")
	var requestBody DeleteUserFinanceCategoryRequestData

	if err := c.ShouldBindBodyWithJSON(&requestBody); err != nil {
		c.JSON(http.StatusOK, common.Fail(err.Error()))
		return
	}

	financeCategory, userErr := handler.categorySvc.FindUserCategoryByUsernameAndCode(username, requestBody.Code)
	if userErr != nil {
		c.JSON(http.StatusOK, common.Fail(userErr.Error()))
		return
	}

	deleteErr := handler.categorySvc.DeleteUserFinanceCategoryAndRelatedRecord(financeCategory)
	if deleteErr != nil {
		c.JSON(http.StatusOK, common.Fail(deleteErr.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Success(""))
}

type DeleteUserFinanceCategoryRequestData struct {
	Code string `json:"code" binding:"required"`
}
