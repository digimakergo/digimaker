//Author xc, Created on 2019-04-03 21:00
//{COPYRIGHTS}

package util

import (
	"os"
	"path/filepath"

	"github.com/digimakergo/digimaker/core/log"

	"github.com/spf13/viper"
)

var appFolder = ""
var configFile = "dm"
var configFolder = "configs"

func AppFolder() string {
	if appFolder == "" {
		//Init App config
		appPath := os.Getenv("dmapp")
		if appPath == "" {
			log.Info("dmapp env not set, use current directory")
			var err error
			appPath, err = os.Getwd()
			if err != nil {
				log.Fatal(err.Error())
			}
			if !FileExists(appPath + "/configs/dm.yaml") {
				log.Fatal("Not a digimaker working directory.")
			}
		}

		if _, err := os.Stat(appPath); os.IsNotExist(err) {
			log.Fatal("Folder " + appPath + " doesn't exist.")
		}

		abs, _ := filepath.Abs(appPath)
		log.Info("Set configurations under " + abs)

		appFolder = appPath
	}
	return appFolder
}

//running mode: eg. dev, prod
func RunningMode() string {
	mode := "prod"
	if env := os.Getenv("env"); env != "" {
		mode = env
	}
	return mode
}

func HomePath() string {
	return AppFolder()
}

func AbsHomePath() string {
	path := HomePath()
	abs, _ := filepath.Abs(path)
	return abs
}

func ConfigPath() string {
	return HomePath() + "/" + configFolder
}

//DMPath returns folder path of the framework. It can be used to load system file(eg. internal setting file)
func DMPath() string {
	return os.Getenv("GOPATH") + "/src/github.com/digimakergo/digimaker"
}

func VarFolder() string {
	varFolder := GetConfig("general", "var_folder")
	result := varFolder
	if !filepath.IsAbs(varFolder) {
		result, _ = filepath.Abs(AbsHomePath() + "/" + varFolder)
	}
	return result
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
		filename = configFile
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

//return nil if section doesn't exist
func GetConfigSectionAll(section string, configFile string) interface{} {
	all := GetAll(configFile)
	value, exist := all[section]
	if !exist {
		return nil
	}
	return value
}

func GetAll(configFile string) map[string]interface{} {
	viper.SetConfigName(configFile)
	viper.AddConfigPath(ConfigPath())
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

//convert viper map interface{} configuration into array map.
func ConvertToMap(config interface{}) map[string]interface{} {
	configMap := config.(map[interface{}]interface{})
	result := map[string]interface{}{}
	for identifier, value := range configMap {
		result[identifier.(string)] = value
	}
	return result
}
