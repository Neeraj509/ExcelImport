package controllers

import (
	"encoding/json"
	"myapp/config"
	"myapp/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// handleError is a helper function for error handling and response
func handleError(c *gin.Context, statusCode int, message string, err error) {
	config.Logger.Printf("%s: %v", message, err)
	c.JSON(statusCode, gin.H{"error": message, "details": err.Error()})
}

// createRecordFromRow creates a Record from an Excel row
func createRecordFromRow(row []string) models.Record {
	return models.Record{
		FirstName:   row[0],
		LastName:    row[1],
		CompanyName: row[2],
		Address:     row[3],
		City:        row[4],
		County:      row[5],
		Postal:      row[6],
		Phone:       row[7],
		Email:       row[8],
		Web:         row[9],
	}
}

// batchInsertAndCache performs batch insertion and caching of records
func batchInsertAndCache(records []models.Record) error {
	if err := config.DB.Create(&records).Error; err != nil {
		return err
	}

	pipe := config.RDB.Pipeline()
	for _, record := range records {
		recordJSON, err := json.Marshal(record)
		if err != nil {
			config.Logger.Printf("Error marshaling record to JSON: %v", err)
			continue
		}
		pipe.Set(config.Ctx, "record:"+strconv.Itoa(int(record.ID)), recordJSON, 5*time.Minute)
	}

	_, err := pipe.Exec(config.Ctx)
	return err
}
