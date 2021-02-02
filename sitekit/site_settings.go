package sitekit

import "github.com/digimakergo/digimaker/core/contenttype"

var siteSettings = map[string]SiteSettings{}
var siteIdentifiers = []string{}

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
	return siteIdentifiers
}

func SetSiteSettings(identifier string, settings SiteSettings) {
	siteSettings[identifier] = settings
	siteIdentifiers = append(siteIdentifiers, identifier)
}
