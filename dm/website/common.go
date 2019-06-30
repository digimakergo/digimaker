package website

var siteSettings = map[string]SiteSettings{}

type SiteSettings struct {
	TemplateBase    string
	TemplateFolders []string
}

func GetSiteSettings(identifier string) SiteSettings {
	return siteSettings[identifier]
}

func InitSiteSettings(identifier string, settings SiteSettings) {
	siteSettings[identifier] = settings
}
