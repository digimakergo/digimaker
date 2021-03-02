package digimaker

import (
	"strconv"

	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/digimakergo/digimaker/rest"
	"github.com/digimakergo/digimaker/sitekit"
	"github.com/gorilla/mux"
)

//Bootstrap digimaker:
//set related path, load definition, load policies
//Route rest and site
func Bootstrap(router *mux.Router) {
	log.Info("Starting from " + util.AbsHomePath())

	router.Use(rest.InitRequest)

	//Route rest api
	restRouter := router.PathPrefix("/api").Subrouter() //todo: make /api configuable
	rest.HandleRoute(restRouter)

	//Route sites
	sitesConfig := util.GetConfigSectionAll("sites", "dm")
	if sitesConfig != nil {
		for i, item := range sitesConfig.([]interface{}) {
			siteConfig := util.ConvertToMap(item)
			err := sitekit.InitSite(router, siteConfig)
			if err != nil {
				log.Error("Error when loading site "+strconv.Itoa(i)+". Detail: "+err.Error(), "")
			}
		}

		//todo: route custom url before content

		//Handle content
		err := sitekit.RouteContent(router)
		if err != nil {
			log.Error("Error when routing sites. Error: "+err.Error(), "")
		}
	}
}

//Initialize db
func InitDB() bool {
	return true
}

func Reload() {

}
