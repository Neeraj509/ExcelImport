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

// ViewRecords retrieves records, first from Redis, then from MySQL
func ViewRecords(c *gin.Context) {
	records, err := getRecordsFromCache()
	if err != nil || len(records) == 0 {
		if err := config.DB.Find(&records).Error; err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to retrieve records from database", err)
			return
		}
		cacheRecords(records)
	}

	c.JSON(http.StatusOK, records)
}

func getRecordsFromCache() ([]models.Record, error) {
	var records []models.Record
	keys, err := config.RDB.Keys(config.Ctx, "record:*").Result()
	if err != nil {
		return records, err
	}

	for _, key := range keys {
		data, err := config.RDB.Get(config.Ctx, key).Result()
		if err != nil {
			config.Logger.Printf("Error retrieving record from Redis with key %s: %v", key, err)
			continue
		}

		var record models.Record
		if err := json.Unmarshal([]byte(data), &record); err != nil {
			config.Logger.Printf("Error unmarshaling record from JSON: %v", err)
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

func cacheRecords(records []models.Record) error {
	pipe := config.RDB.Pipeline()
	for _, record := range records {
		recordJSON, err := json.Marshal(record)
		if err != nil {
			config.Logger.Printf("Error marshaling record to JSON for caching: %v", err)
			continue
		}
		pipe.Set(config.Ctx, "record:"+strconv.Itoa(int(record.ID)), recordJSON, 5*time.Minute)
	}

	if _, err := pipe.Exec(config.Ctx); err != nil {
		return err
	}
	return nil
}
