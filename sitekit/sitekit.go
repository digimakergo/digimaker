package sitekit

import (
	"dm/dm/handler"
	"dm/dm/util"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/flosch/pongo2.v2"
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

		siteSettings := SiteSettings{TemplateBase: templateFolder[0],
			TemplateFolders: templateFolder,
			RootContent:     rootContent,
			DefaultContent:  defaultContent,
			Routes:          routesConfig}
		InitSiteSettings(siteIdentifier, siteSettings)
	}
	return nil
}

func OutputContent(w http.ResponseWriter, r *http.Request, id int, siteIdentifier string, prefix string) {
	// Execute the template per HTTP request
	pongo2.DefaultSet.Debug = true
	siteSettings := GetSiteSettings(siteIdentifier)
	pongo2.DefaultSet.SetBaseDirectory("../templates/" + siteSettings.TemplateBase)
	tpl := pongo2.Must(pongo2.FromCache("../default/content/view.html"))
	querier := handler.Querier()
	content, err := querier.FetchByID(id)
	//todo: handle error, template compiling much better.
	if err != nil {
		fmt.Println(err)
	}
	if !util.ContainsInt(content.GetLocation().Path(), siteSettings.RootContent.GetLocation().ID) {
		w.Write([]byte("Not valid content"))
		return
	}

	err = tpl.ExecuteWriter(pongo2.Context{"content": content,
		"root":     siteSettings.RootContent,
		"viewmode": "full",
		"site":     "demosite",
		"prefix":   prefix}, w)

}

func OutputTemplate(w http.ResponseWriter, r *http.Request, siteIdentifier string, templatePath string, matchData ...map[string]interface{}) {
	realPath := ""
	if len(matchData) == 0 {
		realPath = MatchTemplate(templatePath, map[string]interface{}{})
	} else {
		realPath = MatchTemplate(templatePath, matchData[0])
	}
	if realPath == "" {
		realPath = "default/" + templatePath + ".html"
	}
	siteSettings := GetSiteSettings(siteIdentifier)
	pongo2.DefaultSet.Debug = true
	pongo2.DefaultSet.SetBaseDirectory("../templates/" + siteSettings.TemplateBase)
	fmt.Println(realPath)
	tpl := pongo2.Must(pongo2.FromCache("../" + realPath))
	tpl.Execute(pongo2.Context{})
}
