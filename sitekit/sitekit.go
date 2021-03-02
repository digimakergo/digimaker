package sitekit

import (
	"context"
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

//Initialize all the sites
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

//Handle contents after initialization
func RouteContent(r *mux.Router) error {
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
			ctx := r.Context()

			err := OutputContent(w, id, site, prefix, ctx)
			if err != nil {
				log.Error(err.Error(), "template", r.Context())
				requestID := log.GetContextInfo(ctx).RequestID
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
func OutputContent(w io.Writer, id int, siteIdentifier string, prefix string, ctx context.Context) error {
	siteSettings := GetSiteSettings(siteIdentifier)
	data := map[string]interface{}{
		"root":     siteSettings.RootContent,
		"viewmode": "full",
		"prefix":   prefix}

	querier := handler.Querier()
	content, err := querier.FetchByID(id)
	//todo: handle error, template compiling much better.
	if content == nil {
		data["error"] = "Content not found" //todo: use error code so can we customize it in template
	} else {
		if !util.ContainsInt(content.GetLocation().Path(), siteSettings.RootContent.GetLocation().ID) {
			data["error"] = "Content not found in this site"
		}
	}
	data["content"] = content

	//todo: use anoymouse user id and check permission
	err = Output(w, siteIdentifier, data, ctx)
	return err
}

type TemplateContext struct {
	RequestContext context.Context
	Site           string
}

//Output using template
func Output(w io.Writer, siteIdentifier string, variables map[string]interface{}, ctx context.Context) error {
	// siteSettings := GetSiteSettings(siteIdentifier)

	// pongo2.DefaultSet.SetBaseDirectory("../templates/" + siteSettings.TemplateBase)
	if log.GetContextInfo(ctx).CanDebug() {
		variables["debug"] = true
	}
	gopath := os.Getenv("GOPATH")
	tpl := pongo2.Must(pongo2.FromCache(gopath + "/src/github.com/digimakergo/digimaker/sitekit/templates/main.html")) //todo: use configuration

	variables["site"] = siteIdentifier

	tCtx := TemplateContext{RequestContext: ctx, Site: siteIdentifier}

	for name, newFunctions := range allFunctions {
		functions := newFunctions()
		functions.SetContext(tCtx)
		functionMap := functions.GetMap()
		variables[name] = functionMap
	}

	err := tpl.ExecuteWriter(pongo2.Context(variables), w)
	return err
}
