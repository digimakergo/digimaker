package functions

import (
	"strconv"
	"strings"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/query"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/digimakergo/digimaker/sitekit"
	"github.com/digimakergo/digimaker/sitekit/niceurl"
)

//default dm functions
type DMFunctions struct {
	sitekit.RequestInfo
}

func (dm DMFunctions) GetMap() map[string]interface{} {
	result := map[string]interface{}{
		//if there is site, use site template folder, otherwise use template directly
		"tpl_content": func(content contenttype.ContentTyper, mode string) string {
			siteIdentifer := ""
			//if site is empty, use empty siteIdentifier
			if dm.Site != "" {
				settings := sitekit.GetSiteSettings(dm.Site)
				siteIdentifer = settings.Site
			}

			path := sitekit.GetContentTemplate(dm.Context, content, mode, siteIdentifer)

			log.Debug("Template for content "+content.GetName()+", mode "+mode+": "+path, "template", dm.Context)
			if path != "" && !sitekit.TemplateExist(path) {
				log.Warning("Template file not found: "+path, "template", dm.Context)
			}

			return path
		},

		//Note: site template folder will apply if there is site in matchData
		"tpl_match": func(matchData interface{}, viewType string) string {
			var data map[string]interface{}
			if matchData == nil {
				data = map[string]interface{}{}
			} else {
				data = matchData.(map[string]interface{})
			}

			path, matchLog := sitekit.MatchTemplate(viewType, data)

			log.Debug("Template for view "+viewType+": "+path+"log:"+strings.Join(matchLog, "\n"), "template", dm.Context)
			if path != "" && !sitekit.TemplateExist(path) {
				log.Warning("Template file not found: "+path, "template", dm.Context)
			}

			return path
		},

		"map": func(params ...interface{}) map[string]interface{} {
			result := map[string]interface{}{}
			for i := 0; i < len(params); i = i + 2 {
				key := params[i].(string)
				value := params[i+1]
				result[key] = value
			}
			return result
		},

		"fieldtype": func(field string, content contenttype.ContentTyper) string {
			def, _ := definition.GetDefinition(content.ContentType())
			if fieldDef, ok := def.FieldMap[field]; ok {
				return fieldDef.FieldType
			}
			return ""
		},

		"fetch_byid": func(id int) contenttype.ContentTyper {
			content, err := query.FetchByID(dm.Context, id)
			if err != nil {
				log.Debug("Error when fetch ", "tempalte", dm.Context)
			}
			return content
		},

		"parent": func(content contenttype.ContentTyper) contenttype.ContentTyper {
			parentID := content.Value("parent_id").(int)
			parent, err := query.FetchByID(dm.Context, parentID)
			if err != nil {
				log.Debug("Error when fetch parent", "tempalte", dm.Context)
			}
			return parent
		},

		"children": func(parent contenttype.ContentTyper, contenttype string) []contenttype.ContentTyper {
			userID := util.CurrentUserID(dm.Context)
			children, _, err := query.Children(dm.Context, userID, parent, contenttype, db.EmptyCond().Sortby("l.priority desc", "id asc"))
			if err != nil {
				log.Debug("Error when fetch children on id "+strconv.Itoa(parent.GetID())+": "+err.Error(), "template", dm.Context)
			}
			return children
		},

		"niceurl": func(content contenttype.ContentTyper) string {
			root := sitekit.GetSiteSettings(dm.Site).RootContent
			url := niceurl.GenerateUrl(content, root, dm.SitePath)
			return url
		},

		"root": func(url string) string {
			result := "/" + dm.SitePath
			if url != "/" {
				result = result + "/" + url
			}
			return result
		},
	}

	return result
}

func (dm *DMFunctions) SetInfo(info sitekit.RequestInfo) {
	dm.RequestInfo = info
}

func init() {
	sitekit.RegisterFunctions("dm", func() sitekit.TemplateFunctions {
		return &DMFunctions{}
	})
}
