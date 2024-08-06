package handler

import (
	"GoProject/common"
	"GoProject/model"
	"GoProject/module/category"
	"GoProject/module/user"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type StatisticsHttpHandler struct {
	categorySvc category.Service
	userSvc     user.Service
}

func NewStatisticsHttpHandler(categorySvc category.Service, userSvc user.Service, server *gin.Engine) StatisticsHttpHandler {
	handler := StatisticsHttpHandler{categorySvc: categorySvc, userSvc: userSvc}

	//yearly report
	server.GET("/statistics/:username/yearly", func(c *gin.Context) {
		log.Printf("[%s]%s", c.Request.Method, c.Request.URL)
		handler.getYearlyReport(c)
	})

	//monthly report
	server.GET("/statistics/:username/monthly", func(c *gin.Context) {
		log.Printf("[%s]%s", c.Request.Method, c.Request.URL)
		handler.getMonthReport(c)
	})

	return handler
}

func (handler *StatisticsHttpHandler) getYearlyReport(c *gin.Context) {
	var username = c.Param("username")
	var startYear = c.DefaultQuery("start_year", "")
	var endYear = c.DefaultQuery("end_year", "")

	var startDateTime time.Time
	var endDateTime time.Time
	if startYear != "" {
		startDateTime1, err := time.Parse("2006", startYear)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid date format"})
			return
		}
		startDateTime = startDateTime1
	}
	if endYear != "" {
		startDateTime1, err := time.Parse("2006", endYear)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid date format"})
			return
		}
		endDateTime = startDateTime1.AddDate(1, 0, 0)
	}
	users, userErr := handler.userSvc.FindUser(username)
	if userErr != nil {
		c.JSON(http.StatusOK, common.Fail(userErr.Error()))
		return
	}

	records, recordErr := handler.categorySvc.FindUserRecordsByUserIdPreload(users.ID, startDateTime, endDateTime)
	if recordErr != nil {
		c.JSON(http.StatusOK, common.Fail(recordErr.Error()))
		return
	}

	yearMap := make(map[int][]model.UserFinanceRecord)
	for _, record := range records {
		year := record.SpendDate.Year()
		yearMap[year] = append(yearMap[year], record)
	}

	var result []YearlyReportResponseData
	for key, value := range yearMap {
		detailMap := make(map[string]uint)
		for _, splitRecord := range value {
			detailMap[splitRecord.UserFinanceCategory.Code] += splitRecord.Price
		}

		var resultDetail []ReportDetailResponseData
		for detailCode, detailValue := range detailMap {
			resultDetail = append(resultDetail, ReportDetailResponseData{Category: detailCode, Spend: detailValue})
		}

		result = append(result, YearlyReportResponseData{Year: key, Details: resultDetail})
	}

	c.JSON(http.StatusOK, common.Success(result))
}

type YearlyReportResponseData struct {
	Year    int
	Details []ReportDetailResponseData
}

func (handler *StatisticsHttpHandler) getMonthReport(c *gin.Context) {
	var username = c.Param("username")
	var year = c.DefaultQuery("year", "")
	if year == "" {
		c.JSON(400, gin.H{"error": "No data"})
		return
	}
	var startMonth = year + "-01-01"

	var startDateTime time.Time
	var endDateTime time.Time
	startDateTime1, err := time.Parse("2006-01-02", startMonth)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid date format"})
		return
	}
	startDateTime = startDateTime1.AddDate(0, 0, -1)
	endDateTime = startDateTime1.AddDate(1, 0, -1)

	users, userErr := handler.userSvc.FindUser(username)
	if userErr != nil {
		c.JSON(http.StatusOK, common.Fail(userErr.Error()))
		return
	}

	records, recordErr := handler.categorySvc.FindUserRecordsByUserIdPreload(users.ID, startDateTime, endDateTime)
	if recordErr != nil {
		c.JSON(http.StatusOK, common.Fail(recordErr.Error()))
		return
	}

	yearMap := make(map[string][]model.UserFinanceRecord)
	for _, record := range records {
		yearAndMonth := record.SpendDate.Format("2006-01")
		yearMap[yearAndMonth] = append(yearMap[yearAndMonth], record)
	}

	var result []MonthlyReportResponseData
	for key, value := range yearMap {
		detailMap := make(map[string]uint)
		for _, splitRecord := range value {
			detailMap[splitRecord.UserFinanceCategory.Code] += splitRecord.Price
		}

		var resultDetail []ReportDetailResponseData
		for detailCode, detailValue := range detailMap {
			resultDetail = append(resultDetail, ReportDetailResponseData{Category: detailCode, Spend: detailValue})
		}

		parts := strings.Split(key, "-")
		sendYear, _ := strconv.Atoi(parts[0])
		sendMonth, _ := strconv.Atoi(parts[1])
		result = append(result, MonthlyReportResponseData{
			Year:    sendYear,
			Month:   sendMonth,
			Details: resultDetail,
		})
	}

	c.JSON(http.StatusOK, common.Success(result))
}

type MonthlyReportResponseData struct {
	Year    int
	Month   int
	Details []ReportDetailResponseData
}

type ReportDetailResponseData struct {
	Category string
	Spend    uint
}
