package functions

import (
	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
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
			children, _, err := query.Children(dm.Context, parent, contenttype, userID, db.EmptyCond().Sort("priority desc", "id asc"), false)
			if err != nil {
				log.Debug("Error when fetch ", "tempalte", dm.Context)
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
