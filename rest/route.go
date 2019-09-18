package rest

import (
	"github.com/gorilla/mux"
)

func Route(r *mux.Router) {
	r.HandleFunc("/content/get/{id}", GetContent)
	r.HandleFunc("/content/treemenu/{id}", TreeMenu)
	r.HandleFunc("/util/uploadfile", UploadFile)
	r.HandleFunc("/util/uploadimage", UploadImage)

}
