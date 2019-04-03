//Author xc, Created on 2019-03-27 20:00
//{COPYRIGHTS}

//Package utils provides general utils for the project
package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
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
	fmt.Println("["+level+"]", fmt.Sprint(message...))
}

func GetConfig(section string, identifier string, filename ...string) (string, error) {
	return "mysql", nil
}

func GetConfigSection(section string, configName ...string) map[string]string {
	result := make(map[string]string)
	list := GetConfigSectionI(section, configName...)
	for identifier, value := range list {
		result[identifier] = value.(string)
	}
	return result
}

//Get section of the config,
// config: config file, eg. content(will look for content.yaml or content.json with overriding)
func GetConfigSectionI(section string, configName ...string) map[string]interface{} {
	var filename string
	if configName == nil {
		filename = DefaultSettings.ConfigFile
	} else {
		filename = configName[0]
	}

	viper.SetConfigName(filename)
	viper.AddConfigPath(DefaultSettings.ConfigFolder)
	//todo: support override in section&setting level with order.

	err := viper.ReadInConfig()
	if err != nil {
		Error("Fatal error config file: ", err.Error())
	}
	var result map[string]interface{}
	value := viper.Get(section)
	if value == nil {
		Warning("Section ", section, " doesn't exist on ", filename)
		result = nil
	} else {
		result = value.(map[string]interface{})
	}
	return result
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
