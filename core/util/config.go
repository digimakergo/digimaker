//Author xc, Created on 2019-04-03 21:00
//{COPYRIGHTS}

package util

import (
	"os"
	"path/filepath"

	"github.com/xc/digimaker/core/log"

	"github.com/spf13/viper"
)

var defaultSettings = struct {
	ConfigFile   string
	ConfigFolder string
	HomePath     string
	DMPath       string
}{"dm", "", "", os.Getenv("GOPATH") + "/src/github.com/xc/digimaker"}

func InitHomePath(homePath string) {
	defaultSettings.HomePath = homePath
	defaultSettings.ConfigFolder = defaultSettings.HomePath + "/configs"
}

func HomePath() string {
	return defaultSettings.HomePath
}

func AbsHomePath() string {
	abs, _ := filepath.Abs(defaultSettings.HomePath)
	return abs
}

func ConfigPath() string {
	return defaultSettings.ConfigFolder
}

//DMPath returns folder path of the framework. It can be used to load system file(eg. internal setting file)
func DMPath() string {
	return defaultSettings.DMPath
}

//Get config based on section and identifer
func GetConfig(section string, identifier string, config ...string) string {
	configList := GetConfigSection(section, config...)
	result := configList[identifier]
	return result
}

//Get config string array
func GetConfigArr(section string, identifier string, config ...string) []string {
	configList := GetConfigSectionI(section, config...)
	if _, ok := configList[identifier]; !ok {
		log.Warning("Identifier "+identifier+" doesn't exist on section "+section, "config")
		return nil
	}
	listValue := configList[identifier].([]interface{})
	result := []string{}
	for _, item := range listValue {
		result = append(result, item.(string))
	}
	return result
}

//Get sections with string values
func GetConfigSection(section string, config ...string) map[string]string {
	result := make(map[string]string)
	list := GetConfigSectionI(section, config...)
	for identifier, value := range list {
		result[identifier] = value.(string)
	}
	return result
}

//Get section of the config,
//config: config file, eg. content(will look for content.yaml or content.json with overriding)
func GetConfigSectionI(section string, config ...string) map[string]interface{} {
	var filename string
	if config == nil {
		filename = defaultSettings.ConfigFile
	} else {
		filename = config[0]
	}

	sectionValue := GetConfigSectionAll(section, filename)

	var result map[string]interface{}
	if sectionValue == nil {
		log.Warning("Section "+section+" doesn't exist on "+filename, "")
		result = nil
	} else {
		result = sectionValue.(map[string]interface{})
	}
	return result
}

func GetConfigSectionAll(section string, config string) interface{} {
	all := GetAll(config)
	value, exist := all[section]
	if !exist {
		log.Warning("Key "+section+" doesn't exist in config "+config, "config")
	}
	return value
}

func GetAll(config string) map[string]interface{} {
	viper.SetConfigName(config)
	viper.AddConfigPath(defaultSettings.ConfigFolder)
	//todo: support override in section&setting level with order.
	//todo: did viper cached all? need to verify.
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("Fatal error config file: "+err.Error(), "")
		return nil
	}
	allSettings := viper.AllSettings()
	return allSettings
}

var internalSettings map[string]interface{}

// GetInternalSetting return setting for internal use. This setting is only changed when there is version update.
func GetInternalSetting(setting string) interface{} {
	if value, ok := internalSettings[setting]; ok {
		return value
	} else {
		log.Error("Didn't find setting "+setting+" in dm_internal.yaml", "system")
		return nil
	}
}

func init() {
	viper.SetConfigName("dm_internal")
	viper.AddConfigPath(defaultSettings.DMPath + "/core")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("Fatal error in dm_internal.yaml config file: "+err.Error(), "system")
		return
	}
	internalSettings = viper.AllSettings()
}
