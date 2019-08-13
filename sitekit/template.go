package sitekit

import (
	"dm/core/contenttype"
	"dm/core/util"
	"fmt"
)

const templateViewContent = "content_view"

//Get content view template.
func GetContentTemplate(content contenttype.ContentTyper, viewmode string, settings SiteSettings) string {
	templateRootFolder := util.ConfigPath() + "/../templates"

	matchData := map[string]interface{}{}
	matchData["viewmode"] = viewmode
	matchData["contenttype"] = content.ContentType()
	location := content.GetLocation()
	if location != nil {
		matchData["id"] = location.ID
		matchData["under"] = location.Path()
		matchData["level"] = location.Depth
		matchData["section"] = location.Section
	}

	path := MatchTemplate(templateViewContent, matchData)

	templateFolders := settings.TemplateFolders
	result := ""
	//get the match based on template folder order
	for _, folder := range templateFolders {
		if path != "" {
			pathWithTemplateFolder := folder + "/" + path
			fmt.Println(templateFolders)
			if util.FileExists(templateRootFolder + "/" + pathWithTemplateFolder) {
				result = pathWithTemplateFolder
				break
			}
		}
	}
	return result
}

//MatchTemplate returns overrided template based on override rule in template_override.yaml
func MatchTemplate(source string, matchData map[string]interface{}) string {
	fmt.Println(source)
	rules := util.GetConfigSectionAll(source, "template_override").([]interface{})
	result := ""
	for _, item := range rules {
		conditions := map[string]interface{}{}
		to := ""
		for key, value := range item.(map[interface{}]interface{}) {
			keyStr := key.(string)
			if keyStr == "to" {
				to = value.(string) //todo: have a better name instead of to
				continue
			}
			conditions[keyStr] = value
		}

		matched, matchLog := util.MatchCondition(conditions, matchData)
		if matched {
			result = to
			break
		}
		fmt.Println(matchLog) //todo: add log into somewhere.
	}

	return result
}

//Get template folder list based on site identifier
func TemplateFolders(siteIdentifier string) []string {
	folders := util.GetConfigSectionI("template")["folder"].([]interface{})
	var result = make([]string, len(folders))
	for i, value := range folders {
		result[i] = value.(string)
	}
	return result
}
