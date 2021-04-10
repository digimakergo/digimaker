package functions

import (
	"strconv"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/handler"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/query"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/digimakergo/digimaker/sitekit"
	"github.com/digimakergo/digimaker/sitekit/niceurl"
	"gopkg.in/flosch/pongo2.v2"
)

//default dm functions
type DMFunctions struct {
	sitekit.RequestInfo
}

func (dm DMFunctions) GetMap() map[string]interface{} {
	result := map[string]interface{}{

		"tpl_content": func(content contenttype.ContentTyper, mode string) string {
			settings := sitekit.GetSiteSettings(dm.Site)
			path := sitekit.GetContentTemplate(content, mode, settings, dm.Context)
			log.Debug("Template for content "+content.GetName()+", mode "+mode+": "+path, "template", dm.Context)
			return path
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
		"wash_field": func(input string) *pongo2.Value {
			//wash content attribute
			//todo: support more fieldtypes
			result := handler.ConvertToHtml(dm.Context, input, true, "/var/") //todo: change to use configuration
			return pongo2.AsSafeValue(result)
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
