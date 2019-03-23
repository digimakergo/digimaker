package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

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
	fmt.Printf(message)
}

//UnmarshalData Load json and unmall into variable
func UnmarshalData(filepath string, v interface{}) error {
	file, err := os.Open(filepath)
	if err != nil {
		LogError("Error when loading content definition: " + err.Error())
		return err
	}
	byteValue, _ := ioutil.ReadAll(file)

	err = json.Unmarshal(byteValue, v)
	if err != nil {
		LogError("Error when loading datatype definition: " + err.Error())
		return err
	}

	defer file.Close()
	return nil
}
