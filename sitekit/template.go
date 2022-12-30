package sitekit

import (
	"context"
	"strconv"
	"strings"

	"github.com/digimakergo/digimaker/core/config"
	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/spf13/viper"
)

const templateViewContent = "content_view"
const overrideFile = "template_override"

var overrideFieldtypes = []string{"radio", "select", "checkbox"}

//TemplateFolder() returns folder of templates. eg. under "templates" or "web/templates"
func TemplateFolder() string {
	path := config.AbsHomePath() + "/" + viper.GetString("general.template_folder")
	return path
}

//Get content view template, with site setting's template folder
//todo: support page content(render content) in matchData, eg. "render_field/display_type": "news", "render_contentype": "frontpage" - used like 'section' isolation based on visit page.
func GetContentTemplate(ctx context.Context, content contenttype.ContentTyper, viewmode string, siteIdentifier string) string {

	matchData := map[string]interface{}{}
	matchData["viewmode"] = viewmode
	matchData["site"] = siteIdentifier
	matchData["contenttype"] = content.ContentType()
	location := content.GetLocation()
	if location != nil {
		matchData["id"] = location.ID
		matchData["under"] = location.Path()
		matchData["level"] = location.Depth
		matchData["section"] = location.Section
	}
	for field, fieldDef := range content.Definition().FieldMap {
		if util.Contains(overrideFieldtypes, fieldDef.FieldType) {
			matchData["field/"+field] = content.Value(field)
		}
	}

	matchLog := []string{}
	path, matchLog := MatchTemplate(ctx, templateViewContent, matchData)

	log.Debug("Matching on "+content.GetName()+", got: "+path+"\n "+strings.Join(matchLog, "\n"), "template-match", ctx)
	return path
}

//MatchTemplate returns overrided template based on override config(eg. template_override.yaml)
func MatchTemplate(ctx context.Context, viewSection string, matchData map[string]interface{}, fileName ...string) (string, []string) {
	overrideFileName := ""
	result := ""
	matchLog := []string{}
	if len(fileName) == 0 {
		overrideFileName = overrideFile
		//if there is include, match in included file
		viper := config.GetViper(overrideFileName)
		includeI := viper.Get("include")
		if includeI != nil {
			for _, item := range includeI.([]interface{}) {
				includeRules := map[string]interface{}{}
				includedFile := ""
				templateFolder := ""
				for key, value := range item.(map[interface{}]interface{}) {
					keyS := key.(string)
					if keyS == "file" {
						includedFile = value.(string)
					} else if keyS == "template_folder" {
						templateFolder = value.(string)
					} else {
						includeRules[keyS] = value
					}
				}
				includeMatched, _ := util.MatchCondition(includeRules, matchData)
				if includeMatched {
					matchLog = append(matchLog, "Matching on include file: "+includedFile)
					var includedMatchLog []string
					result, includedMatchLog = MatchTemplate(ctx, viewSection, matchData, includedFile)
					matchLog = append(matchLog, includedMatchLog...)
					if result != "" {
						result = WashTemplatePath(result, templateFolder)
						return result, matchLog
					} else {
						matchLog = append(matchLog, "Not matched on include file : "+includedFile)
					}
				}
			}
		}
	} else {
		overrideFileName = fileName[0]
	}

	overrideViper := config.GetViper(overrideFileName)
	rulesI := overrideViper.Get(viewSection)
	if rulesI == nil {
		return "", []string{"view section not found: " + viewSection}
	}
	rules := rulesI.([]interface{})
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

		matchLog = append(matchLog, "Matching on rule"+strconv.Itoa(i)+" on file "+overrideFile)
		matched, currentMatchLog := util.MatchCondition(conditions, matchData)
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

	result = WashTemplatePath(result, "")
	return result, matchLog
}
