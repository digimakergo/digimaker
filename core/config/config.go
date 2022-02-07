//Author xc, Created on 2019-04-03 21:00
//{COPYRIGHTS}

package config

import (
	"os"
	"path/filepath"

	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/util"

	"github.com/spf13/viper"
)

var appFolder = ""
var configFile = "dm"
var configFolder = "configs"

const envPrefix = "DM"

func initDefaultViper() {
	log.Info("Initializing config")
	viper.AutomaticEnv()
	viper.SetConfigName(configFile)
	viper.SetEnvPrefix(envPrefix)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(ConfigPath())
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("Fatal error config file: "+err.Error(), "")
	}
}

func HomePath() string {
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
			if !util.FileExists(appPath + "/configs/dm.yaml") {
				log.Fatal("Not a digimaker working directory.")
			}
		}

		if _, err := os.Stat(appPath); os.IsNotExist(err) {
			log.Fatal("Folder " + appPath + " doesn't exist.")
		}

		abs, _ := filepath.Abs(appPath)
		log.Info("Set configurations under " + abs)

		appFolder = appPath

		initDefaultViper()
	}
	return appFolder
}

func AbsHomePath() string {
	path := HomePath()
	abs, _ := filepath.Abs(path)
	return abs
}

func ConfigPath() string {
	return HomePath() + "/" + configFolder
}

func VarFolder() string {
	varFolder := viper.GetString("general.var_folder")
	result := varFolder
	if !filepath.IsAbs(varFolder) {
		result, _ = filepath.Abs(AbsHomePath() + "/" + varFolder)
	}
	return result
}

var viperMap map[string]*viper.Viper = map[string]*viper.Viper{}

func GetViper(configFile string) *viper.Viper {
	_, ok := viperMap[configFile]
	if !ok {
		configFile = util.SecurePath(configFile)
		v := viper.New()
		v.SetConfigName(configFile)
		v.SetConfigType("yaml")
		v.AddConfigPath(ConfigPath())
		err := v.ReadInConfig()
		if err != nil {
			log.Error("Fatal error config file: "+err.Error(), "")
			return nil
		}
		viperMap[configFile] = v
	}

	result := viperMap[configFile]
	return result
}
