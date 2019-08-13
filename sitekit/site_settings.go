package sitekit

import "dm/core/contenttype"

var siteSettings = map[string]SiteSettings{}

//a basic setting to run a site.
type SiteSettings struct {
	TemplateBase    string
	TemplateFolders []string
	RootContent     contenttype.ContentTyper
	DefaultContent  contenttype.ContentTyper
	Routes          []interface{} //host, path.
}

func GetSiteSettings(identifier string) SiteSettings {
	return siteSettings[identifier]
}

func GetSites() []string {
	var result []string
	for key, _ := range siteSettings {
		result = append(result, key)
	}
	return result
}

func SetSiteSettings(identifier string, settings SiteSettings) {
	siteSettings[identifier] = settings
}
