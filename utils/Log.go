package utils

import (
	"log"
	"os"
)

const error = "error"
const warning = "warning"
const notice = "notice"

func LogError(message string) {
	Log(message, "error")
}

func LogWarning(message string) {
	Log(message, "warning")
}

func Notice(message string) {
	Log(message, "notice")
}

/*
Log message
*/
func Log(message string, level string) {
	logDirectory := "log"
	path := logDirectory + level

	logMessage := "[<ip>]" + message

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)
	log.Printf(logMessage)
}
