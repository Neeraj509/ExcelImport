package controllers

import (
	"encoding/json"
	"myapp/config"
	"myapp/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// EditRecord edits a record by ID
func EditRecord(c *gin.Context) {
	id := c.Param("id")

	var record models.Record
	if err := config.DB.Where("id = ?", id).First(&record).Error; err != nil {
		handleError(c, http.StatusNotFound, "Record not found", err)
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid input", err)
		return
	}

	if err := config.DB.Model(&record).Updates(updates).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to update record", err)
		return
	}

	cacheRecord(record)
	c.JSON(http.StatusOK, record)
}

func cacheRecord(record models.Record) error {
	recordJSON, err := json.Marshal(record)
	if err != nil {
		config.Logger.Printf("Error marshaling updated record to JSON: %v", err)
		return err
	}

	if err := config.RDB.Set(config.Ctx, "record:"+strconv.Itoa(int(record.ID)), recordJSON, 5*time.Minute).Err(); err != nil {
		config.Logger.Printf("Error setting record in Redis: %v", err)
		return err
	}

	return nil
}
