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
const overrideFile = "template_override"

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
	matchData["site"] = settings.Site
	matchData["contenttype"] = content.ContentType()
	location := content.GetLocation()
	if location != nil {
		matchData["id"] = location.ID
		matchData["under"] = location.Path()
		matchData["level"] = location.Depth
		matchData["section"] = location.Section
	}

	path := ""
	matchLog := []string{}

	siteOverride := overrideFile + "-" + settings.Site
	if util.FileExists(util.ConfigPath() + "/" + siteOverride + ".yaml") { //todo: use viper way so json/tomal can be supported also
		currentPath, currentMatchLog := MatchTemplate(templateViewContent, matchData, siteOverride)
		path = currentPath
		matchLog = append(matchLog, currentMatchLog...)
	}
	if path == "" {
		currentPath, currentMatchLog := MatchTemplate(templateViewContent, matchData, overrideFile)
		path = currentPath
		matchLog = append(matchLog, currentMatchLog...)
	}

	log.Debug("Matching on "+content.GetName()+", got: "+path+"\n "+strings.Join(matchLog, "\n"), "template-match", ctx)

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

//MatchTemplate returns overrided template based on override config(eg. template_override.yaml)
func MatchTemplate(viewSection string, matchData map[string]interface{}, fileName ...string) (string, []string) {
	overrideFileName := overrideFile
	if len(fileName) > 0 {
		overrideFileName = fileName[0]
	}
	rulesI := util.GetConfigSectionAll(viewSection, overrideFileName)
	if rulesI == nil {
		return "", []string{"view section not found: " + viewSection}
	}
	rules := rulesI.([]interface{})
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
		matchLog = append(matchLog, "matching on rule"+strconv.Itoa(i)+" on file "+overrideFile)
		matchLog = append(matchLog, currentMatchLog...)
		if matched {
			washedVars := map[string]string{}
			for key, value := range matchData {
				switch value.(type) {
				case string:
					washedVars[key] = value.(string)
				case int:
					washedVars[key] = strconv.Itoa(value.(int))
				}
			}
			result = util.ReplaceStrVar(to, washedVars)
			break
		}
	}

	return result, matchLog
}
