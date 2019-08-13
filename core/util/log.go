//Author xc, Created on 2019-05-04 22:28
//{COPYRIGHTS}
//This is a debug based on context.

package util

import (
	"fmt"
	"log"
	"os"
)

var logger *log.Logger

func Logger() *log.Logger {
	if logger == nil {
		logger = log.New(os.Stdout, "", log.LstdFlags) //todo: use file and separate error and warning in file.
	}
	return logger
}

func Debug(message ...interface{}) {
	debugEnabled := GetConfig("debug", "enabled") == "yes"
	if debugEnabled {
		Log("debug", message...)
	}
}

func Error(message ...interface{}) {
	Log("error", message...)
}

func Warning(message ...interface{}) {
	Log("warning", message...)
}

/*
Log message
*/
func Log(logType string, message ...interface{}) {
	/*
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
	*/
	//todo: log into files
	//todo: for into client screen in debug mode.
	//todo: include cagetory as parameter
	logger := Logger()
	logger.Println("["+logType+"]", fmt.Sprint(message...))
}
