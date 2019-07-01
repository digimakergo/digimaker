package sitekit

import "dm/dm/contenttype"

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

func InitSiteSettings(identifier string, settings SiteSettings) {
	siteSettings[identifier] = settings
}
