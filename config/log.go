package config

import (
	"log"
	"os"
)

// Logger is the custom logger instance
var Logger *log.Logger

// InitLogger initializes the logger and writes to a log file
func InitLogger() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Create a new logger instance
	Logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}
