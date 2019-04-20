//Author xc, Created on 2019-03-27 20:00
//{COPYRIGHTS}

//Package utils provides general utils for the project
package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/rs/xid"
)

func Error(message ...interface{}) {
	Log("error", message...)
}

func Warning(message ...interface{}) {
	Log("warning", message...)
}

func Notice(message ...interface{}) {
	Log("notice", message...)
}

func Debug(category string, message ...interface{}) {
	Log("debug,"+category, message...)
}

/*
Log message
*/
func Log(level string, message ...interface{}) {
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
	fmt.Println("["+level+"]", fmt.Sprint(message...))
}

//Add time point for calcuate how much time takes for operations.
//Typical type include: database, operation, template
//Typical identifier include: layout.tpl when comes to template, add when it comes to operation
func AddTimePoint(category string, identifier string) {

}

//UnmarshalData Load json and unmall into variable
func UnmarshalData(filepath string, v interface{}) error {
	file, err := os.Open(filepath)
	if err != nil {
		Error("Error when loading content definition: " + err.Error())
		return err
	}
	byteValue, _ := ioutil.ReadAll(file)

	err = json.Unmarshal(byteValue, v)
	if err != nil {
		Error("Error when loading datatype definition: " + err.Error())
		return err
	}

	defer file.Close()
	return nil
}

//Generate unique id. It should be cluster safe.
func GenerateUID() string {
	guid := xid.New()
	guidStr := guid.String()
	return guidStr
}

//Convert like "hello_world" to "HelloWorld"
func UpperName(input string) string {
	arr := strings.Split(input, "_")
	for i := range arr {
		arr[i] = strings.Title(arr[i])
	}
	return strings.Join(arr, "")
}
