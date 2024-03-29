package functions

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/query"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/digimakergo/digimaker/core/util/image"
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
			def, _ := definition.GetDefinition(content.ContentType())
			if fieldDef, ok := def.FieldMap[field]; ok {
				return fieldDef.FieldType
			}
			return ""
		},

		"fetch_byid": func(id int) contenttype.ContentTyper {
			content, err := query.FetchByLID(dm.Context, id)
			if err != nil {
				log.Debug("Error when fetch ", "tempalte", dm.Context)
			}
			return content
		},

		"fetch_by_cid": func(ctype string, cid int) contenttype.ContentTyper {
			content, err := query.FetchByCID(dm.Context, ctype, cid)
			if err != nil {
				log.Debug("Error when fetch ", "tempalte", dm.Context)
			}
			return content
		},

		"parent": func(content contenttype.ContentTyper) contenttype.ContentTyper {
			parentID := content.Value("parent_id").(int)
			parent, err := query.FetchByLID(dm.Context, parentID)
			if err != nil {
				log.Debug("Error when fetch parent", "tempalte", dm.Context)
			}
			return parent
		},

		"sublist": func(parent contenttype.ContentTyper, contenttype string, depth int, params ...interface{}) []contenttype.ContentTyper {
			userID := util.CurrentUserID(dm.Context)
			def, _ := definition.GetDefinition(contenttype)
			sortBy := []string{}
			if def.HasLocation {
				sortBy = []string{"l.priority desc", "id asc"}
			} else {
				sortBy = []string{"published desc"}
			}
			cond := paramsToCondition(db.EmptyCond().Sortby(sortBy...), params...)
			list, _, err := query.SubList(dm.Context, userID, parent, contenttype, depth, cond)
			if err != nil {
				log.Debug("Error when fetch children on id "+strconv.Itoa(parent.GetID())+": "+err.Error(), "template", dm.Context)
			}
			return list
		},

		"sublist_count": func(parent contenttype.ContentTyper, contenttype string, depth int, cond db.Condition) int {
			userID := util.CurrentUserID(dm.Context)
			_, count, err := query.SubList(dm.Context, userID, parent, contenttype, depth, cond.Limit(0, 0))
			if err != nil {
				log.Debug("Error when fetch children on id "+strconv.Itoa(parent.GetID())+": "+err.Error(), "template", dm.Context)
			}
			return count
		},

		"children": func(parent contenttype.ContentTyper, contenttype string, params ...interface{}) []contenttype.ContentTyper {
			userID := util.CurrentUserID(dm.Context)
			var sortBy []string
			def, _ := definition.GetDefinition(contenttype)
			if def.HasLocation {
				sortBy = []string{"l.priority desc", "id asc"}
			} else {
				sortBy = []string{"published desc"}
			}

			cond := db.EmptyCond().Sortby(sortBy...)
			cond = paramsToCondition(cond, params...)

			children, _, err := query.Children(dm.Context, userID, parent, contenttype, cond)
			if err != nil {
				log.Debug("Error when fetch children on id "+strconv.Itoa(parent.GetID())+": "+err.Error(), "template", dm.Context)
			}
			return children
		},

		"children_count": func(parent contenttype.ContentTyper, contenttype string, params ...interface{}) int {
			userID := util.CurrentUserID(dm.Context)

			cond := db.EmptyCond()
			if len(params) > 0 {
				if param, ok := params[0].(db.Condition); ok {
					cond = param
				}
			}

			_, count, err := query.Children(dm.Context, userID, parent, contenttype, cond.Limit(0, 0).WithCount())
			if err != nil {
				log.Debug("Error when fetch children on id "+strconv.Itoa(parent.GetID())+": "+err.Error(), "template", dm.Context)
			}
			return count
		},

		"niceurl": func(content contenttype.ContentTyper) string {
			root := sitekit.GetSiteSettings(dm.Site).RootContent
			url := niceurl.GenerateUrl(content, root, dm.SitePath)
			return url
		},

		"output_field": func(field string, content contenttype.ContentTyper) interface{} {
			outputValue, err := query.OutputField(dm.Context, content, field)
			if err != nil {
				log.Error(fmt.Errorf("Error when output_field %v: %v", field, err.Error()), "", dm.Context)
				return ""
			}
			return outputValue
		},

		//convert image path to real
		//params: size - "original", "default", "800", ...
		"image": func(path string, size ...string) string {
			sizeStr := "original"
			if len(size) > 0 {
				sizeStr = size[0]
			}
			result := image.ImagePath(dm.Context, path, sizeStr)
			return result
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

		"now": func() time.Time {
			return time.Now()
		},
	}

	return result
}

/**
Convert parameters to condition, with sort offset, etc
0 - sortby string: priority desc, publish desc
1 - limit, eg. 10
2 - condition. eg db.Cond("author", 10)
3 - offset
*/
func paramsToCondition(cond db.Condition, params ...interface{}) db.Condition {
	sortBy := []string{}

	if len(params) == 0 {
		return cond
	}

	limit := -1
	offset := 0
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

	if len(params) >= 3 {
		if param2, ok := params[2].(db.Condition); ok {
			cond = cond.And(param2)
		}
	}

	if len(params) >= 4 {
		offset = params[3].(int)
	}
	result := cond.Sortby(sortBy...).Limit(offset, limit)
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
