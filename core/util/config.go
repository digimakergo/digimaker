//Author xc, Created on 2019-04-03 21:00
//{COPYRIGHTS}

package util

import (
	"dm/core/log"
	"os"

	"github.com/spf13/viper"
)

var defaultSettings = struct {
	ConfigFile   string
	ConfigFolder string
	HomePath     string
	PackageName  string
}{"site", "", "", ""}

func SetPackageName(packageName string) {
	defaultSettings.PackageName = packageName
	defaultSettings.HomePath = os.Getenv("GOPATH") + "/src/" + packageName //todo: change to not using gopath
	defaultSettings.ConfigFolder = defaultSettings.HomePath + "/configs"
}

func HomePath() string {
	return defaultSettings.HomePath
}

func PackageName() string {
	return defaultSettings.PackageName
}

func ConfigPath() string {
	return defaultSettings.ConfigFolder
}

//Get config based on section and identifer
func GetConfig(section string, identifier string, config ...string) string {
	configList := GetConfigSection(section, config...)
	result := configList[identifier]
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
	viper.SetConfigName(config)
	viper.AddConfigPath(defaultSettings.ConfigFolder)
	//todo: support override in section&setting level with order.
	//todo: did viper cached all? need to verify.
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("Fatal error config file: "+err.Error(), "")
		return nil
	}

	value := viper.Get(section)
	return value
}
