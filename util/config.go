//Author xc, Created on 2019-04-03 21:00
//{COPYRIGHTS}

package util

import (
	"github.com/spf13/viper"
)

//DefaultSettingsType and DefaultSettings set all the default,
// including which config folder it is.
type DefaultSettingsType struct {
	ConfigFile   string
	ConfigFolder string
}

var defaultSettings = DefaultSettingsType{"site", "/Users/xc/go/caf-prototype/src/dm/configs"}

func SetConfigPath(folderPath string) {
	defaultSettings.ConfigFolder = folderPath
}

func ConfigPath() string {
	return defaultSettings.ConfigFolder
}

//Get config based on section and identifer
func GetConfig(section string, identifier string, config ...string) string {
	configList := GetConfigSection(section, config...)
	result := configList[section]
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

	viper.SetConfigName(filename)
	viper.AddConfigPath(defaultSettings.ConfigFolder)
	//todo: support override in section&setting level with order.
	//todo: did viper cached all? need to verify.

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

func GetAll(config string) {

}
