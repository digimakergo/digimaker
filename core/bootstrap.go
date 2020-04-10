package core

import (
	"github.com/xc/digimaker/core/contenttype"
	"github.com/xc/digimaker/core/fieldtype"
	"github.com/xc/digimaker/core/permission"
	"github.com/xc/digimaker/core/util"
	"github.com/xc/digimaker/core/log"
)

func Bootstrap(packageName string) bool {
	log.Info("Starting from " + packageName)
	util.SetPackageName(packageName)
	err := contenttype.LoadDefinition()
	if err != nil {
		return false
	}
	err = fieldtype.LoadDefinition()
	if err != nil {
		return false
	}
	err = permission.LoadPolicies()
	if err != nil {
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
