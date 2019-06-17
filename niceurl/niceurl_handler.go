//Package niceurl provides nice url feature for dm framework
package niceurl

import (
	"dm/dm/contenttype"
	"strconv"
)

func GenerateUrl(content contenttype.ContentTyper) string {
	location := content.GetLocation()
	result := ""
	if location != nil {
		path := location.IdentifierPath
		pattern := "digit" //todo: read from config file.
		switch pattern {
		case "digit":
			result = path + "-" + strconv.Itoa(location.ID)
		default:

		}
	} else {
		//todo: give a warning.
	}
	return result
}
