package controllers

import (
	"myapp/config"
	"myapp/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteRecord deletes a record by ID
func DeleteRecord(c *gin.Context) {
	id := c.Param("id")

	// Attempt to delete the record from the database
	if err := config.DB.Where("id = ?", id).Delete(&models.Record{}).Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to delete record from database", err)
		return
	}

	// Attempt to delete the record from Redis cache
	cacheKey := "record:" + id
	if err := config.RDB.Del(config.Ctx, cacheKey).Err(); err != nil {
		config.Logger.Printf("Error deleting record from Redis cache with key %s: %v", cacheKey, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}

// DeleteAllRecords deletes all records from the database and Redis
func DeleteAllRecords(c *gin.Context) {
	// Attempt to delete all records from the database
	if err := config.DB.Exec("DELETE FROM records").Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to delete records from the database", err)
		return
	}

	// Reset the auto-increment value (for MySQL)
	if err := config.DB.Exec("ALTER TABLE records AUTO_INCREMENT = 1").Error; err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to reset auto-increment value", err)
		return
	}
	// Attempt to retrieve all Redis keys for deletion
	keys, err := config.RDB.Keys(config.Ctx, "record:*").Result()
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to retrieve keys from Redis", err)
		return
	}
	// Use a Redis pipeline for deleting multiple keys
	pipe := config.RDB.Pipeline()
	for _, key := range keys {
		_, err := pipe.Del(config.Ctx, key).Result()
		if err != nil {
			config.Logger.Printf("Error deleting key %s from Redis: %v", key, err)
		}
	}

	// Execute the pipeline
	if _, err := pipe.Exec(config.Ctx); err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to execute Redis pipeline", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All records deleted successfully"})
}
