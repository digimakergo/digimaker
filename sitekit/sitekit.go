package sitekit

import (
	"context"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/permission"
	"github.com/digimakergo/digimaker/core/query"
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
	Site            string
	TemplateBase    string
	TemplateFolders []string
	RootContent     contenttype.ContentTyper
	DefaultContent  contenttype.ContentTyper
	Host            string
	Path            string
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

//Initialize sites setting to memory
func LoadSite(siteConfig map[string]interface{}) error {
	siteIdentifier := siteConfig["identifier"].(string)

	if _, ok := siteConfig["template_folder"]; !ok {
		return errors.New("Need template_folder setting.")
	}
	templateFolder := util.InterfaceToStringArray(siteConfig["template_folder"].([]interface{}))

	if _, ok := siteConfig["root"]; !ok {
		return errors.New("Need root setting.")
	}
	root := siteConfig["root"].(int)
	rootContent, err := query.FetchByID(context.Background(), root)
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
		defaultContent, err = query.FetchByID(context.Background(), defaultInt)
		if err != nil {
			return errors.Wrap(err, "Default doesn't exist.")
		}
	}

	host := siteConfig["host"].(string)
	path := ""
	if _, ok := siteConfig["path"]; ok {
		path = siteConfig["path"].(string)
	}
	siteSettings := SiteSettings{
		Site:            siteIdentifier,
		TemplateBase:    templateFolder[0],
		TemplateFolders: templateFolder,
		RootContent:     rootContent,
		DefaultContent:  defaultContent,
		Host:            host,
		Path:            path}
	SetSiteSettings(siteIdentifier, siteSettings)
	log.Info("Site settings loaded: " + siteIdentifier)
	return nil
}

//Handle content, given mux variables: site: <site>, path: <site path>, id: <content id>
func HandleContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	sitePath := GetSitePath(r)
	site := vars["site"]

	ctx := r.Context()
	err := OutputContent(w, id, site, sitePath, ctx)
	if err != nil {
		log.Error(err.Error(), "template", r.Context())
		requestID := log.GetContextInfo(ctx).RequestID
		http.Error(w, "Error occurred. request id: "+requestID, http.StatusInternalServerError)
	}
}

// remove / if path has / in the end
func GetSitePath(r *http.Request) string {
	sitePath := ""
	vars := mux.Vars(r)
	if path, ok := vars["path"]; ok {
		sitePath = path
	}
	result := strings.TrimSuffix(sitePath, "/")
	return result
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	site := vars["site"]

	defaultContent := siteSettings[site].DefaultContent
	defaultContentID := defaultContent.GetLocation().ID
	vars["id"] = strconv.Itoa(defaultContentID)

	r = mux.SetURLVars(r, vars)
	HandleContent(w, r)
}

func setVar(r *http.Request, key string, value string) *http.Request {
	vars := mux.Vars(r)
	vars[key] = value
	r = mux.SetURLVars(r, vars)
	return r
}

type SiteRouters map[string]*mux.Router

func GetSiteRouters(r *mux.Router) (SiteRouters, SiteRouters) {
	//loop sites and route
	sites := GetSites()

	subRouters := SiteRouters{}
	defaultRouters := SiteRouters{}
	for _, identifier := range sites {
		settings := GetSiteSettings(identifier)
		host := settings.Host
		path := settings.Path

		var subRouter *mux.Router
		//use subrouter which is better for performance
		if path != "" {
			subRouter = r.Host(host).PathPrefix("/{path:" + path + "}/").Subrouter()

			defaultRouter := r.Host(host).PathPrefix("/{path:" + path + "}").Subrouter()
			defaultRouters[identifier] = defaultRouter
		} else {
			subRouter = r.Host(host).Subrouter()
			defaultRouters[identifier] = subRouter
		}
		subRouters[identifier] = subRouter
	}
	return subRouters, defaultRouters
}

//Handle contents after initialization
func RouteContent(siteRouters SiteRouters, defaultRouters SiteRouters) {
	for site, subRouter := range siteRouters {
		subRouter.HandleFunc("/content/view/{id}", func(w http.ResponseWriter, r *http.Request) {
			r = setVar(r, "site", mux.CurrentRoute(r).GetName())
			HandleContent(w, r)
		}).Name(site)

		subRouter.MatcherFunc(niceurl.ViewContentMatcher).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = setVar(r, "site", mux.CurrentRoute(r).GetName())
			HandleContent(w, r)
		}).Name(site)
	}

	for site, router := range defaultRouters {
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			r = setVar(r, "site", mux.CurrentRoute(r).GetName())
			handleRoot(w, r)
		}).Name(site)

		router.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
			r = setVar(r, "site", mux.CurrentRoute(r).GetName())
			handleRoot(w, r)
		}).Name(site)
	}
}

//Output content using conent template
func OutputContent(w io.Writer, id int, siteIdentifier string, sitePath string, ctx context.Context) error {
	siteSettings := GetSiteSettings(siteIdentifier)
	variables := map[string]interface{}{
		"root":     siteSettings.RootContent,
		"default":  siteSettings.DefaultContent,
		"viewmode": "full",
		"site":     siteIdentifier,
		"sitepath": sitePath}

	content, err := query.FetchByID(ctx, id)

	if err != nil {
		return errors.Wrap(err, "Error of outputing content while fetching content")
	}

	if content == nil {
		variables["error"] = "Content not found" //todo: use error code and set variables(from template) so can we customize it in template
	} else {
		if !util.ContainsInt(content.GetLocation().Path(), siteSettings.RootContent.GetLocation().ID) {
			variables["error"] = "Content not found in this site"
		}

		userID := util.CurrentUserID(ctx)
		if !permission.CanRead(ctx, userID, content) {
			variables["error"] = "No permission to this content"
		}
	}
	variables["content"] = content

	err = Output(w, variables, ctx)
	return err
}

type RequestInfo struct {
	Context  context.Context
	Site     string
	SitePath string
}

//Output using template
func Output(w io.Writer, variables map[string]interface{}, ctx context.Context) error {
	// siteSettings := GetSiteSettings(siteIdentifier)

	// pongo2.DefaultSet.SetBaseDirectory("../templates/" + siteSettings.TemplateBase)
	if log.GetContextInfo(ctx).CanDebug() {
		variables["debug"] = true
	}
	gopath := os.Getenv("GOPATH")
	tpl := pongo2.Must(pongo2.FromCache(gopath + "/src/github.com/digimakergo/digimaker/sitekit/templates/main.html")) //todo: use configuration

	info := RequestInfo{Context: ctx, Site: variables["site"].(string), SitePath: variables["sitepath"].(string)}

	for name, newFunctions := range allFunctions {
		functions := newFunctions()
		functions.SetInfo(info)
		functionMap := functions.GetMap()
		variables[name] = functionMap
	}

	err := tpl.ExecuteWriter(pongo2.Context(variables), w)
	return err
}

func init() {
	//Route sites
	sitesConfig := util.GetConfigSectionAll("sites", "dm")
	if sitesConfig != nil {
		for i, item := range sitesConfig.([]interface{}) {
			siteConfig := util.ConvertToMap(item)
			err := LoadSite(siteConfig)
			if err != nil {
				log.Fatal("Error when loading site " + strconv.Itoa(i) + ". Detail: " + err.Error())
			}
		}
	}
}
