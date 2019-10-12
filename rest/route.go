package rest

import (
	"github.com/gorilla/mux"
)

func Route(r *mux.Router) {
	r.HandleFunc("/contenttype/get/{contentype}", GetDefinition)
	r.HandleFunc("/content/get/{id}", GetContent)
	r.HandleFunc("/content/treemenu/{id}", TreeMenu)
	r.HandleFunc("/content/list/{id}", Children)

	r.HandleFunc("/content/new/{parent}/{contenttype}", New)

	r.HandleFunc("/form/validate/{contenttype}", Validate)

	r.HandleFunc("/util/uploadfile", UploadFile)
	r.HandleFunc("/util/uploadimage", UploadImage)

}
