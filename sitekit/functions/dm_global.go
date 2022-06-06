package functions

import (
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/sitekit"
)

//default dm functions
type DMGlobalFunctions struct {
}

func (dm DMGlobalFunctions) GetMap() map[string]interface{} {
	result := map[string]interface{}{
		"cond": func(field string, value interface{}) db.Condition {
			return db.Cond(field, value)
		},
	}

	return result
}

func (dm *DMGlobalFunctions) SetInfo(info sitekit.RequestInfo) {
}

func init() {
	sitekit.RegisterFunctions("_", func() sitekit.TemplateFunctions {
		return &DMGlobalFunctions{}
	})
}
