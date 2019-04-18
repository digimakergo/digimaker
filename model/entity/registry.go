package entity

import "dm/model"

func NewInstance(contentType string) interface{} {
	// var result model.ContentTyper
	// switch contentType {
	// case "article":
	// 	result = Article{}
	// case "folder":
	// 	result = Folder{}
	// }
	return Article{}
}

//Global variable for registering contentType
//todo: support collection, eg.[]article to bind collection easier from db.
var contenttypeRegistry = map[string]func() model.ContentTyper{}

func Register(contentType string, newContentType func() model.ContentTyper) {
	contenttypeRegistry[contentType] = newContentType
}

func NewContentType(contentType string) model.ContentTyper {
	return contenttypeRegistry[contentType]()
}

func init() {
	Register("article", func() model.ContentTyper {
		return Article{}
	})
	Register("folder", func() model.ContentTyper {
		return Folder{}
	})
}
