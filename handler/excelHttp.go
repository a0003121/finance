package handler

import (
	"GoProject/common"
	"GoProject/model"
	"GoProject/module/category"
	"GoProject/module/user"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ExcelHttpHandler struct {
	categorySvc category.Service
	userSvc     user.Service
}

func NewExcelHttpHandler(userSvc user.Service, categorySvc category.Service, server *gin.Engine) ExcelHttpHandler {
	var handler = ExcelHttpHandler{userSvc: userSvc, categorySvc: categorySvc}

	server.GET("/user/:username/excel", func(c *gin.Context) {
		log.Printf("[%s]%s", c.Request.Method, c.Request.URL)
		//_ = jwt.Authenticate(c) && jwt.IsAdmin(c)
	}, func(c *gin.Context) {
		handler.exportExcelRecords(c)
	})

	server.POST("/user/:username/excel_upload", func(c *gin.Context) {
		log.Printf("[%s]%s", c.Request.Method, c.Request.URL)
		//_ = jwt.Authenticate(c) && jwt.IsAdmin(c)
	}, func(c *gin.Context) {
		handler.excelImport(c)
	})

	server.GET("/user/:username/excel_template", func(c *gin.Context) {
		log.Printf("[%s]%s", c.Request.Method, c.Request.URL)
		//_ = jwt.Authenticate(c) && jwt.IsAdmin(c)
	}, func(c *gin.Context) {
		handler.exportExcelImportTemplate(c)
	})

	return handler
}

func (h ExcelHttpHandler) excelImport(c *gin.Context) {
	file, getFileErr := c.FormFile("file")
	if getFileErr != nil {
		c.JSON(http.StatusOK, common.Fail(getFileErr.Error()))
		return
	}

	// Open the file
	src, openFileErr := file.Open()
	if openFileErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open the file"})
		return
	}
	defer src.Close()

	// Create a temporary file to store the uploaded file
	tempFile, tempFileErr := os.CreateTemp("", "upload-*.xlsx")
	if tempFileErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a temporary file"})
		return
	}
	defer os.Remove(tempFile.Name())

	// Copy the uploaded file to the temporary file
	if _, copyFileErr := io.Copy(tempFile, src); copyFileErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the uploaded file"})
		return
	}

	// Open the temporary file with excelize
	f, openFileErr := excelize.OpenFile(tempFile.Name())
	if openFileErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read the Excel file"})
		return
	}

	var username = c.Param("username")
	users, userErr := h.userSvc.FindUser(username)
	if userErr != nil {
		c.JSON(http.StatusOK, common.Fail(userErr.Error()))
		return
	}
	userCategory, categoryErr := h.categorySvc.FindUserCategoriesByUsername(username)
	if categoryErr != nil {
		c.JSON(http.StatusOK, common.Fail(categoryErr.Error()))
	}
	categoryMap := make(map[string]uint, len(userCategory))
	for _, userCategory := range userCategory {
		categoryMap[userCategory.Code] = userCategory.ID
	}

	sheetName := f.GetSheetName(1)
	rows := f.GetRows(sheetName)
	var userRecords []model.UserFinanceRecord
	for index, row := range rows {
		if index == 0 {
			continue
		}
		if row[0] != "" && row[1] != "" && row[2] != "" {
			var spendDate = row[0]
			var inputCategory = row[1]
			var price = row[2]
			if _, ok := categoryMap[inputCategory]; !ok {
				c.JSON(http.StatusOK, common.Fail("input category not exist"))
				return
			}

			dateTime, dateErr := time.Parse("01-02-06", spendDate)
			if dateErr != nil {
				c.JSON(http.StatusOK, common.Fail("Invalid date format"))
				return
			}

			num, numErr := strconv.ParseUint(price, 10, 32)
			if numErr != nil {
				c.JSON(http.StatusOK, common.Fail("Invalid num format"))
				return
			}

			userRecords = append(userRecords, model.UserFinanceRecord{
				UsersID:               users.ID,
				UserFinanceCategoryId: categoryMap[inputCategory],
				Price:                 uint(num),
				SpendDate:             dateTime,
			})
		} else {
			break
		}
	}

	if len(userRecords) == 0 {
		c.JSON(http.StatusOK, common.Fail("no data to import"))
		return
	}

	createErr := h.categorySvc.CreateUserFinanceRecords(&userRecords)
	if createErr != nil {
		c.JSON(http.StatusOK, common.Fail(createErr.Error()))
		return
	}

	// Respond with the cell value
	c.JSON(http.StatusOK, common.Success(""))
}

func (h ExcelHttpHandler) exportExcelImportTemplate(c *gin.Context) {
	var username = c.Param("username")

	userCategory, categoryErr := h.categorySvc.FindUserCategoriesByUsername(username)
	if categoryErr != nil {
		c.JSON(http.StatusOK, common.Fail(categoryErr.Error()))
	}

	var categories []string
	for _, financeCategory := range userCategory {
		categories = append(categories, financeCategory.Code)
	}

	file := excelize.NewFile()
	sheetName := "Sheet1"

	file.SetCellValue(sheetName, "A1", "日期")
	file.SetCellValue(sheetName, "B1", "類別")
	file.SetCellValue(sheetName, "C1", "金額")
	file.SetColWidth(sheetName, "A", "C", 20)

	dvRange := excelize.NewDataValidation(true)
	dvRange.Sqref = "B:B"
	dvRange.SetDropList(categories)
	file.AddDataValidation(sheetName, dvRange)

	dateStyle, styleErr := file.NewStyle(`{
		"number_format": 14
	}`)
	if styleErr != nil {
		c.JSON(http.StatusOK, common.Fail(styleErr.Error()))
	}

	numberStyle, styleErr := file.NewStyle(`{
		"number_format": 1
	}`)
	if styleErr != nil {
		c.JSON(http.StatusOK, common.Fail(styleErr.Error()))
	}

	file.SetCellStyle(sheetName, "A2", "A50", dateStyle)
	file.SetCellStyle(sheetName, "C2", "C50", numberStyle)

	filePath := "template.xlsx"
	if err := file.SaveAs(filePath); err != nil {
		c.JSON(http.StatusOK, common.Fail(err.Error()))
	}

	c.Header("Content-Disposition", "attachment; filename="+filePath)
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
	defer os.Remove(filePath)
}

func (h ExcelHttpHandler) exportExcelRecords(c *gin.Context) {
	var username = c.Param("username")

	records, recordErr := h.categorySvc.FindUserRecordsByUsernamePreload(username)
	if recordErr != nil {
		c.JSON(http.StatusOK, common.Fail(recordErr.Error()))
	}

	userCategory, categoryErr := h.categorySvc.FindUserCategoriesByUsername(username)
	if categoryErr != nil {
		c.JSON(http.StatusOK, common.Fail(categoryErr.Error()))
	}

	var categories []string
	for _, financeCategory := range userCategory {
		categories = append(categories, financeCategory.Code)
	}

	file := excelize.NewFile()
	sheetName := "Sheet1"

	file.SetCellValue(sheetName, "A1", "日期")
	file.SetCellValue(sheetName, "B1", "類別")
	file.SetCellValue(sheetName, "C1", "金額")
	file.SetColWidth(sheetName, "A", "C", 20)

	for i, record := range records {
		row := i + 2
		file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), record.SpendDate)
		file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), record.UserFinanceCategory.Code)
		file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), record.Price)
	}

	filePath := "record.xlsx"
	if err := file.SaveAs(filePath); err != nil {
		c.JSON(http.StatusOK, common.Fail(err.Error()))
	}

	c.Header("Content-Disposition", "attachment; filename="+filePath)
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
	defer os.Remove(filePath)
}
