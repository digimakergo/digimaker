package sitekit

type TemplateFunctions interface {
	GetMap() map[string]interface{}
	SetContext(ctx TemplateContext)
}

type NewFunctions = func() TemplateFunctions

var allFunctions map[string]NewFunctions

func RegisterFunctions(name string, implementation func() TemplateFunctions) {
	if allFunctions == nil {
		allFunctions = map[string]NewFunctions{}
	}
	allFunctions[name] = implementation
}
