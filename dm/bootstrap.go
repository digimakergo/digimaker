package dm

import (
	"dm/dm/contenttype"
	"dm/dm/fieldtype"
	"dm/dm/permission"
	"dm/dm/util"
)

func Boot(projectHome string) bool {
	util.SetConfigPath(projectHome + "/configs")
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
