package functions

import (
	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/sitekit"
)

//default functions
type DMFunctions struct {
	context sitekit.TemplateContext
}

func (dm DMFunctions) GetMap() map[string]interface{} {
	result := map[string]interface{}{}
	result["tpl_content"] = func(content contenttype.ContentTyper, mode string) string {
		ctx := dm.context
		settings := sitekit.GetSiteSettings(ctx.Site)
		path := sitekit.GetContentTemplate(content, mode, settings)
		log.Debug("Template for content "+content.GetName()+", mode "+mode+": "+path, "template", ctx.RequestContext)
		return path
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
