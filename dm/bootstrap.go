package dm

import (
	"dm/dm/contenttype"
	"dm/dm/fieldtype"
	"dm/dm/permission"
	"dm/dm/util"
)

type Bootstrap struct{}

func (Bootstrap) Boot(home string) bool {
	util.SetConfigPath(home + "/configs")
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

func (Bootstrap) Reload() {

}
