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
	"strconv"
	"time"
)

type RecordHttpHandler struct {
	categorySvc category.Service
	userSvc     user.Service
}

func NewRecordHandler(categorySvc category.Service, userSvc user.Service, server *gin.Engine) RecordHttpHandler {
	var handler = RecordHttpHandler{categorySvc: categorySvc, userSvc: userSvc}

	//get user record
	server.GET("/user/:username/record", func(c *gin.Context) {
		log.Printf("[%s]%s", c.Request.Method, c.Request.URL)
		//_ = jwt.Authenticate(c) && jwt.IsAdmin(c)
	}, func(c *gin.Context) {
		handler.findUserRecords(c)
	})

	//create user record
	server.POST("/user/record", func(c *gin.Context) {
		log.Printf("[%s]%s", c.Request.Method, c.Request.URL)
		//_ = jwt.Authenticate(c) && jwt.IsAdmin(c)
	}, func(c *gin.Context) {
		handler.createUserRecord(c)
	})

	//delete user record
	server.DELETE("/user/record/:recordId", func(c *gin.Context) {
		log.Printf("[%s]%s", c.Request.Method, c.Request.URL)
		//_ = jwt.Authenticate(c) && jwt.IsAdmin(c)
	}, func(c *gin.Context) {
		handler.deleteUserRecord(c)
	})

	//update user record
	server.PUT("/user/:username/record/:recordId", func(c *gin.Context) {
		log.Printf("[%s]%s", c.Request.Method, c.Request.URL)
		//_ = jwt.Authenticate(c) && jwt.IsAdmin(c)
	}, func(c *gin.Context) {
		handler.dpdateUserRecord(c)
	})
	return handler
}

func (handler *RecordHttpHandler) deleteUserRecord(c *gin.Context) {
	var recordId = c.Param("recordId")

	recordIdInt, recordIdErr := strconv.Atoi(recordId) // convert string to int
	if recordIdErr != nil {
		c.JSON(http.StatusOK, common.Fail(recordIdErr.Error()))
		return
	}
	deleteErr := handler.categorySvc.DeleteUserFinanceRecordById(uint(recordIdInt))
	if deleteErr != nil {
		c.JSON(http.StatusOK, common.Fail(deleteErr.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Success(""))
}

func (handler *RecordHttpHandler) dpdateUserRecord(c *gin.Context) {
	var recordId = c.Param("recordId")
	var username = c.Param("username")
	recordIdInt, recordIdErr := strconv.Atoi(recordId) // convert string to int
	if recordIdErr != nil {
		c.JSON(http.StatusOK, common.Fail(recordIdErr.Error()))
		return
	}
	var requestBody UpdateUserRecordRequestData

	if err := c.ShouldBindBodyWithJSON(&requestBody); err != nil {
		c.JSON(http.StatusOK, common.Fail(err.Error()))
		return
	}

	userCategory, categoryErr := handler.categorySvc.FindUserCategoryByUsernameAndCode(username, requestBody.Code)
	if categoryErr != nil {
		c.JSON(http.StatusOK, common.Fail(categoryErr.Error()))
	}

	// Parse the date string to time.Time
	date, err := time.Parse("2006-01-02", requestBody.SpendDate)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid date format"})
		return
	}

	datas := map[string]interface{}{
		"UserFinanceCategoryId": userCategory.ID,
		"Price":                 requestBody.Price,
		"SpendDate":             date,
	}
	updateErr := handler.categorySvc.ModifyUserFinanceRecordById(uint(recordIdInt), datas)
	if updateErr != nil {
		c.JSON(http.StatusOK, common.Fail(updateErr.Error()))
	}

	c.JSON(http.StatusOK, common.Success(""))
}

type UpdateUserRecordRequestData struct {
	Code      string `json:"code" binding:"required"`
	SpendDate string `json:"spend_date" binding:"required"`
	Price     uint   `json:"price" binding:"required"`
}

func (handler *RecordHttpHandler) findUserCategories(c *gin.Context) {
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

type UserCategoryResponseData struct {
	Code       string    `json:"code"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func (handler *RecordHttpHandler) createUserRecord(c *gin.Context) {
	var requestBody CreateUserRecordRequestData

	if err := c.ShouldBindBodyWithJSON(&requestBody); err != nil {
		c.JSON(http.StatusOK, common.Fail(err.Error()))
		return
	}

	// Parse the date string to time.Time
	date, err := time.Parse("2006-01-02", requestBody.SpendDate)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid date format"})
		return
	}

	users, userErr := handler.userSvc.FindUser(requestBody.Username)
	if userErr != nil {
		c.JSON(http.StatusOK, common.Fail(userErr.Error()))
		return
	}

	userCategory, findErr := handler.categorySvc.FindUserCategoryByUsernameAndCode(requestBody.Username, requestBody.Code)
	if findErr != nil {
		c.JSON(http.StatusOK, common.Fail(findErr.Error()))
		return
	}

	var createData = model.UserFinanceRecord{
		UsersID:               users.ID,
		UserFinanceCategoryId: userCategory.ID,
		Price:                 requestBody.Price,
		SpendDate:             date,
	}

	createErr := handler.categorySvc.CreateUserFinanceRecord(&createData)
	if createErr != nil {
		c.JSON(http.StatusOK, common.Fail(createErr.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Success(""))
}

type CreateUserRecordRequestData struct {
	Code      string `json:"code" binding:"required"`
	Username  string `json:"username" binding:"required"`
	SpendDate string `json:"spend_date" binding:"required"`
	Price     uint   `json:"price" binding:"required"`
}

func (handler *RecordHttpHandler) findUserRecords(c *gin.Context) {
	var username = c.Param("username")

	pageNumber := c.DefaultQuery("page_number", "1")
	pageSize := c.DefaultQuery("page_size", "5")
	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")

	var startDateTime time.Time
	var endDateTime time.Time
	if startDate != "" {
		startDateTime1, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid date format"})
			return
		}
		startDateTime = startDateTime1
	}

	if endDate != "" {
		endDateTime2, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid date format"})
			return
		}
		endDateTime = endDateTime2
	}

	pageNumberInt, pageNumberErr := strconv.Atoi(pageNumber) // convert string to int
	if pageNumberErr != nil {
		c.JSON(http.StatusOK, common.Fail(pageNumberErr.Error()))
		return
	}
	pageSizeInt, pageSizeErr := strconv.Atoi(pageSize) // convert string to int
	if pageSizeErr != nil {
		c.JSON(http.StatusOK, common.Fail(pageSizeErr.Error()))
		return
	}
	users, userErr := handler.userSvc.FindUser(username)
	if userErr != nil {
		c.JSON(http.StatusOK, common.Fail(userErr.Error()))
		return
	}

	count, records, recordErr := handler.categorySvc.FindUserRecordsPageByUserIdPreload(users.ID, pageNumberInt, pageSizeInt, startDateTime, endDateTime)
	if recordErr != nil {
		c.JSON(http.StatusOK, common.Fail(recordErr.Error()))
		return
	}

	var result []UserRecordDetailResponseData
	for _, userFinanceRecord := range records {
		result = append(result, UserRecordDetailResponseData{
			ID:        userFinanceRecord.ID,
			Code:      userFinanceRecord.UserFinanceCategory.Code,
			Price:     userFinanceRecord.Price,
			SpendTime: userFinanceRecord.SpendDate.Format("2006-01-02"),
		})
	}

	final := UserRecordResponseData{
		CurrentPage: pageNumberInt,
		PageSize:    pageSizeInt,
		TotalPage:   int(count)/pageSizeInt + 1,
		Details:     result,
	}

	c.JSON(http.StatusOK, common.Success(final))
}

type UserRecordResponseData struct {
	TotalPage   int                            `json:"total_page"`
	CurrentPage int                            `json:"current_page"`
	PageSize    int                            `json:"page_size"`
	Details     []UserRecordDetailResponseData `json:"details"`
}

type UserRecordDetailResponseData struct {
	ID        uint   `json:"id"`
	Code      string `json:"code"`
	Price     uint   `json:"price"`
	SpendTime string `json:"spend_time"`
}
