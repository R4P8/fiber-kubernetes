package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var AppLogger *log.Logger

func InitLogger() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %s", err)
	}

	AppLogger = log.New(logFile, "APP_LOG ", log.Ldate|log.Ltime|log.Lshortfile)
	log.Println("Logger initialized")
}

func LogToFileAndES(level, message string) {
	// Log ke file (dibaca Filebeat)
	AppLogger.Printf("[%s] %s\n", level, message)

	// Kirim ke Elasticsearch
	if ES != nil {
		doc := fmt.Sprintf(`{"level":"%s","message":"%s","timestamp":"%s"}`, level, message, time.Now().Format(time.RFC3339))
		_, err := ES.Index("app-logs", strings.NewReader(doc))
		if err != nil {
			AppLogger.Printf("[ERROR] failed to index log to ES: %v", err)
		}
	}
}
