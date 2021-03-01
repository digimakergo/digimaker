package functions

import (
	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/handler"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/sitekit"
)

//default dm functions
type DMFunctions struct {
	context sitekit.TemplateContext
}

func (dm DMFunctions) GetMap() map[string]interface{} {
	result := map[string]interface{}{

		"tpl_content": func(content contenttype.ContentTyper, mode string) string {
			ctx := dm.context
			settings := sitekit.GetSiteSettings(ctx.Site)
			path := sitekit.GetContentTemplate(content, mode, settings, ctx.RequestContext)
			log.Debug("Template for content "+content.GetName()+", mode "+mode+": "+path, "template", ctx.RequestContext)
			return path
		},

		"fetch_byid": func(id int) contenttype.ContentTyper {
			querier := handler.Querier()
			content, err := querier.FetchByID(id)
			if err != nil {
				log.Debug("Error when fetch ", "tempalte", dm.context.RequestContext)
			}
			return content
		},

		"parent": func(content contenttype.ContentTyper) contenttype.ContentTyper {
			querier := handler.Querier()
			parentID := content.Value("parent_id").(int)
			parent, err := querier.FetchByID(parentID)
			if err != nil {
				log.Debug("Error when fetch parent", "tempalte", dm.context.RequestContext)
			}
			return parent
		},

		"children": func(parent contenttype.ContentTyper, contenttype string) []contenttype.ContentTyper {
			querier := handler.Querier()
			children, _, err := querier.Children(parent, contenttype, 1, db.EmptyCond(), nil, nil, false, dm.context.RequestContext)
			if err != nil {
				log.Debug("Error when fetch ", "tempalte", dm.context.RequestContext)
			}
			return children
		},
	}

	return result
}

func (dm *DMFunctions) SetContext(ctx sitekit.TemplateContext) {
	dm.context = ctx
}

func init() {
	sitekit.RegisterFunctions("dm", func() sitekit.TemplateFunctions {
		return &DMFunctions{}
	})
}
