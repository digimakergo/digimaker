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

			path, matchLog := sitekit.MatchTemplate(dm.Context, viewType, data)

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
			def, _ := definition.GetDefinition(dm.Context, content.ContentType())
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

		"children": func(parent contenttype.ContentTyper, contenttype string, params ...interface{}) []contenttype.ContentTyper {
			userID := util.CurrentUserID(dm.Context)
			limit := -1
			var sortBy []string
			def, _ := definition.GetDefinition(dm.Context, contenttype)
			if def.HasLocation {
				sortBy = []string{"l.priority desc", "published desc"}
			} else {
				sortBy = []string{"published desc"}
			}
			if len(params) > 0 {
				//sort

				if params[0] != "" {
					sortBy = strings.Split(params[0].(string), ",")
				}

				if len(params) >= 2 {
					//0 means no limit(default max limit)
					paramLimit := params[1].(int)
					if paramLimit > 0 {
						limit = paramLimit
					}
				}

			}
			children, _, err := query.Children(dm.Context, userID, parent, contenttype, db.EmptyCond().Sortby(sortBy...).Limit(0, limit))
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

		"request_id": func() string {
			return log.GetContextInfo(dm.Context).RequestID
		},

		"abs_path": func(path string) string {
			//get path under current package
			path = sitekit.TemplateFolder() + "/" + util.SecurePath(path)
			return path
		},

		"output_field": func(content contenttype.ContentTyper, fieldIdentifier string) interface{} {
			value := query.OutputField(dm.Context, content, fieldIdentifier)
			return value
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
