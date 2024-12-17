package logger

import (
	"log"
)

// Info logs informational messages
func Info(message string, args ...interface{}) {
	log.Printf("[INFO] "+message, args...)
}

// Error logs error messages
func Error(message string, args ...interface{}) {
	log.Printf("[ERROR] "+message, args...)
}
