package sitekit

import (
	"context"
	"strconv"
	"strings"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/util"
)

const templateViewContent = "content_view"

//TemplateFolder() returns folder of templates. eg. under "templates" or "web/templates"
func TemplateFolder() string {
	path := util.AbsHomePath() + "/" + util.GetConfig("general", "template_folder")
	return path
}

//Get content view template.
func GetContentTemplate(content contenttype.ContentTyper, viewmode string, settings SiteSettings, ctx context.Context) string {
	templateFolder := TemplateFolder()

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

	path, matchLog := MatchTemplate(templateViewContent, matchData)
	log.Debug("Matching on "+content.GetName()+", got: "+path, "template-match", ctx)
	log.Debug(strings.Join(matchLog, "\n"), "template-match", ctx)

	templateFolders := settings.TemplateFolders
	result := ""
	//get the match based on template folder order
	for _, folder := range templateFolders {
		if path != "" {
			pathWithTemplateFolder := folder + "/" + path
			if util.FileExists(templateFolder + "/" + pathWithTemplateFolder) {
				result = pathWithTemplateFolder
				break
			} else {
				log.Warning("Matched file not found: "+path, "template", ctx)
			}
		}
	}
	return result
}

//MatchTemplate returns overrided template based on override rule in template_override.yaml
func MatchTemplate(source string, matchData map[string]interface{}) (string, []string) {
	rules := util.GetConfigSectionAll(source, "template_override").([]interface{})
	result := ""
	matchLog := []string{}
	for i, item := range rules {
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

		matched, currentMatchLog := util.MatchCondition(conditions, matchData)
		matchLog = append(matchLog, "matching on rule"+strconv.Itoa(i))
		matchLog = append(matchLog, currentMatchLog...)
		if matched {
			result = to
			break
		}
	}

	return result, matchLog
}
