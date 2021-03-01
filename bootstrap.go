package digimaker

import (
	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/permission"
	"github.com/digimakergo/digimaker/core/util"
)

//Bootstrap digimaker:
//set related path, load definition, load policies
func Bootstrap() {
	log.Info("Starting from " + util.AbsHomePath())

	err := contenttype.LoadDefinition()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = permission.LoadPolicies()
	if err != nil {
		log.Fatal("Loading policies error: " + err.Error())
	}
}

//Initialize db
func InitDB() bool {
	return true
}

func Reload() {

}
