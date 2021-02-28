package sitekit

//TemplateFunctions represents a set of functions
//To register a new set of function, first implement this interface and then invoke RegisterFunctions in init()
type TemplateFunctions interface {
	GetMap() map[string]interface{}
	SetContext(ctx TemplateContext)
}

var allFunctions map[string]NewFunctions

type NewFunctions = func() TemplateFunctions

//Register a set of functions
//name: 'namespace' of functions
func RegisterFunctions(name string, implementation func() TemplateFunctions) {
	if allFunctions == nil {
		allFunctions = map[string]NewFunctions{}
	}
	allFunctions[name] = implementation
}
