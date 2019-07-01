package sitekit

import (
	"dm/dm/contenttype"
	"dm/dm/handler"
	"dm/dm/util"
	"dm/niceurl"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"

	"github.com/gorilla/mux"
	"gopkg.in/flosch/pongo2.v2"
)

func InitSites(r *mux.Router, config map[string]interface{}) error {
	for siteIdentifier, item := range config {
		siteConfig := item.(map[string]interface{})

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
			return errors.Wrap(err, "Root doesn't exist.")
		}

		//todo: default can be optional.
		if _, ok := siteConfig["default"]; !ok {
			return errors.New("Need default setting.")
		}
		defaultInt := siteConfig["default"].(int)
		var defaultContent contenttype.ContentTyper
		if defaultInt == root {
			defaultContent = rootContent
		} else {
			defaultContent, err = handler.Querier().FetchByID(defaultInt)
			if err != nil {
				return errors.Wrap(err, "Default doesn't exist.")
			}
		}

		routesConfig := siteConfig["routes"].([]interface{})
		siteSettings := SiteSettings{TemplateBase: templateFolder[0],
			TemplateFolders: templateFolder,
			RootContent:     rootContent,
			DefaultContent:  defaultContent,
			Routes:          routesConfig}
		SetSiteSettings(siteIdentifier, siteSettings)
	}
	return nil
}

func HandleContent(r *mux.Router) error {
	//loop sites and route
	sites := GetSites()
	for _, identifier := range sites {
		var handleContentView = func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			id, _ := strconv.Atoi(vars["id"])
			prefix := ""
			if path, ok := vars["path"]; ok {
				prefix = path
			}
			OutputContent(w, r, id, identifier, prefix)
		}

		//site route and get sub route
		err := SiteRouter(r, identifier, func(s *mux.Router) {
			s.HandleFunc("/content/view/{id}", handleContentView)
			s.MatcherFunc(niceurl.ViewContentMatcher).HandlerFunc(handleContentView)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

//Output content using conent template
func OutputContent(w http.ResponseWriter, r *http.Request, id int, siteIdentifier string, prefix string) {
	querier := handler.Querier()
	content, err := querier.FetchByID(id)
	//todo: handle error, template compiling much better.
	if err != nil {
		fmt.Println(err)
	}

	siteSettings := GetSiteSettings(siteIdentifier)
	if !util.ContainsInt(content.GetLocation().Path(), siteSettings.RootContent.GetLocation().ID) {
		w.Write([]byte("Not valid content"))
		return
	}

	data := map[string]interface{}{"content": content,
		"root":     siteSettings.RootContent,
		"viewmode": "full",
		"prefix":   prefix}
	Output(w, r, siteIdentifier, "content/view", data)
}

//Output using template
func Output(w http.ResponseWriter, r *http.Request, siteIdentifier string, templatePath string, variables map[string]interface{}, matchedData ...map[string]interface{}) {
	siteSettings := GetSiteSettings(siteIdentifier)
	pongo2.DefaultSet.Debug = true
	pongo2.DefaultSet.SetBaseDirectory("../templates/" + siteSettings.TemplateBase)
	tpl := pongo2.Must(pongo2.FromCache("../main.html"))

	variables["site"] = siteIdentifier

	variables["template"] = templatePath
	if len(matchedData) == 0 {
		variables["matched_data"] = nil
	} else {
		variables["matched_data"] = matchedData[0]
	}
	err := tpl.ExecuteWriter(pongo2.Context(variables), w)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
