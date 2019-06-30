package sitekit

import (
	"dm/dm/handler"
	"dm/dm/util"
	"dm/dm/website"
	"errors"

	"github.com/gorilla/mux"
)

func Init(r *mux.Router, config map[string]interface{}) error {
	for siteIdentifier, item := range config {
		siteConfig := item.(map[string]interface{})
		//init site settings

		routesConfig := siteConfig["routes"].([]interface{})

		if _, ok := siteConfig["template_folder"]; !ok {
			return errors.New("Need template_folder setting.")
		}
		templateFolder := util.InterfaceToStringArray(siteConfig["template_folder"].([]interface{}))

		if _, ok := siteConfig["root"]; !ok {
			return errors.New("Need root setting.")
		}
		root := siteConfig["root"].(int)
		rootContent, err := handler.Querier().FetchByID(root)
		if err != nil {
			return err
		}

		if _, ok := siteConfig["default"]; !ok {
			return errors.New("Need default setting.")
		}
		defaultInt := siteConfig["default"].(int)
		defaultContent, err := handler.Querier().FetchByID(defaultInt)
		if err != nil {
			return err
		}

		siteSettings := website.SiteSettings{TemplateBase: templateFolder[0],
			TemplateFolders: templateFolder,
			RootContent:     rootContent,
			DefaultContent:  defaultContent,
			Routes:          routesConfig}
		website.InitSiteSettings(siteIdentifier, siteSettings)
	}
	return nil
}
