package digimaker

import (
	"os"
	"path/filepath"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/permission"
	"github.com/digimakergo/digimaker/core/util"
)

//Bootstrap digimaker:
//set related path, load definition, load policies
func Bootstrap(homePath string) bool {

	if _, err := os.Stat(homePath); os.IsNotExist(err) {
		log.Fatal("Folder " + homePath + " doesn't exist.")
		return false
	}

	abs, _ := filepath.Abs(homePath)
	log.Info("Starting from " + abs)

	util.InitHomePath(homePath)
	err := contenttype.LoadDefinition()
	if err != nil {
		log.Fatal(err.Error())
		return false
	}

	err = permission.LoadPolicies()
	if err != nil {
		log.Fatal("Loading policies error: " + err.Error())
		return false
	}
	return true
}

//Initialize db
func InitDB() bool {
	return true
}

func Reload() {

}
