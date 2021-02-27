package sitekit

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/handler"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/digimakergo/digimaker/sitekit/niceurl"

	"github.com/pkg/errors"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
)

func InitSite(r *mux.Router, siteConfig map[string]interface{}) error {

	siteIdentifier := siteConfig["identifier"].(string)

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
	return nil
}

func HandleContent(r *mux.Router) error {
	//loop sites and route
	sites := GetSites()
	for _, identifier := range sites {
		var handleContentView = func(w http.ResponseWriter, r *http.Request, site string) {
			vars := mux.Vars(r)
			id, _ := strconv.Atoi(vars["id"])
			prefix := ""
			if path, ok := vars["path"]; ok {
				prefix = path
			}
			err := OutputContent(w, id, site, prefix)
			if err != nil {
				log.Error(err.Error(), "template", r.Context())
				requestID := log.GetContextInfo(r.Context()).RequestID
				http.Error(w, "Error occurred. request id: "+requestID, http.StatusInternalServerError)
			}
		}

		//site route and get sub route
		err := SiteRouter(r, identifier, func(s *mux.Router, site string) {
			s.HandleFunc("/content/view/{id}", func(w http.ResponseWriter, r *http.Request) {
				handleContentView(w, r, site)
			})
			s.MatcherFunc(niceurl.ViewContentMatcher).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handleContentView(w, r, site)
			})

			//default page to same as handling content/view/<default>
			s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				defaultContent := siteSettings[identifier].DefaultContent
				defaultContentID := defaultContent.GetLocation().ID
				r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(defaultContentID)})
				handleContentView(w, r, site)
			})
		})
		if err != nil {
			return err
		}

	}
	return nil
}

//Output content using conent template
func OutputContent(w io.Writer, id int, siteIdentifier string, prefix string) error {
	querier := handler.Querier()
	content, err := querier.FetchByID(id)
	//todo: handle error, template compiling much better.
	if content == nil {
		return errors.New("Content not found")
	}

	siteSettings := GetSiteSettings(siteIdentifier)
	if !util.ContainsInt(content.GetLocation().Path(), siteSettings.RootContent.GetLocation().ID) {
		return errors.New("Content not availebl under this site")
	}

	data := map[string]interface{}{"content": content,
		"root":     siteSettings.RootContent,
		"viewmode": "full",
		"prefix":   prefix}
	err = Output(w, siteIdentifier, "content/view", data)
	return err
}

//Output using template
func Output(w io.Writer, siteIdentifier string, templatePath string, variables map[string]interface{}, matchedData ...map[string]interface{}) error {
	// siteSettings := GetSiteSettings(siteIdentifier)
	// pongo2.DefaultSet.Debug = true
	// pongo2.DefaultSet.SetBaseDirectory("../templates/" + siteSettings.TemplateBase)
	gopath := os.Getenv("GOPATH")
	tpl := pongo2.Must(pongo2.FromCache(gopath + "/src/github.com/digimakergo/digimaker/sitekit/templates/main.html"))

	variables["site"] = siteIdentifier

	variables["template"] = templatePath
	if len(matchedData) == 0 {
		variables["matched_data"] = nil
	} else {
		variables["matched_data"] = matchedData[0]
	}
	err := tpl.ExecuteWriter(pongo2.Context(variables), w)
	if err != nil {
		return errors.New(templatePath + ": " + err.Error())
	}
	return nil
}
