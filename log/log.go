package log

import (
	"log"
)

const (
	DebugLevel = iota
	InfoLevel
	ErrorLevel
)

var Level = InfoLevel

func Debug(msg string, args ...interface{}) {
	if Level <= DebugLevel {
		log.Printf("[DEBUG] "+msg, args...)
	}
}

func Info(msg string, args ...interface{}) {
	if Level <= InfoLevel {
		log.Printf("[INFO] "+msg, args...)
	}
}

func Error(msg string, args ...interface{}) {
	if Level <= ErrorLevel {
		log.Printf("[ERROR] "+msg, args...)
	}
}

func Fatal(err error) {
	if Level <= ErrorLevel {
		log.Fatalf("[FATAL] " + err.Error())
	}
}
