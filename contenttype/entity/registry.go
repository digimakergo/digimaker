package entity

import "dm/contenttype"

func NewInstance(contentType string) interface{} {
	var result interface{}
	switch contentType {
	case "article":
		result = &[]Article{}
	case "folder":
		result = &[]Folder{}
	}
	return result
}

//Global variable for registering contentType
//todo: support collection, eg.[]article to bind collection easier from db.
var contenttypeRegistry = map[string]func() contenttype.ContentTyper{}

func Register(contentType string, newContentType func() contenttype.ContentTyper) {
	contenttypeRegistry[contentType] = newContentType
}

func NewContentType(contentType string) contenttype.ContentTyper {
	return contenttypeRegistry[contentType]()
}

func init() {
	Register("article", func() contenttype.ContentTyper {
		return Article{}
	})
	Register("folder", func() contenttype.ContentTyper {
		return Folder{}
	})
}
