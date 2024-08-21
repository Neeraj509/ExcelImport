package controllers

import (
	"myapp/config"
	"myapp/models"
	"net/http"
	"path/filepath"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
)

// ImportExcel handles the Excel upload
func ImportExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		handleError(c, http.StatusBadRequest, "File upload failed", err)
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".xlsx" && ext != ".xls" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only .xlsx and .xls files are allowed"})
		return
	}

	localFilePath := filepath.Join("uploads", file.Filename)
	if err := c.SaveUploadedFile(file, localFilePath); err != nil {
		handleError(c, http.StatusInternalServerError, "Unable to save file", err)
		return
	}

	go processExcel(localFilePath)

	c.JSON(http.StatusOK, gin.H{"message": "File processing started"})
}

func processExcel(filePath string) {

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		config.Logger.Printf("Error opening Excel file: %v", err)
		return
	}

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		config.Logger.Println("No sheets found in the Excel file")
		return
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		config.Logger.Printf("Error getting rows from sheet: %v", err)
		return
	}

	var records []models.Record
	for i, row := range rows {
		if i == 0 {
			continue
		}
		record := createRecordFromRow(row)
		records = append(records, record)
	}

	if err := batchInsertAndCache(records); err != nil {
		config.Logger.Printf("Error inserting records: %v", err)
	}
}
