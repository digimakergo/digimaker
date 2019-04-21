package entity

import "dm/contenttype"

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

func NewList(contentType string) interface{} {
	var result interface{}
	switch contentType {
	case "article":
		result = &[]Article{}
	case "folder":
		result = &[]Folder{}
	}
	return result
}

func NewInstance(contentType string) contenttype.ContentTyper {
	var result contenttype.ContentTyper
	switch contentType {
	case "article":
		result = &Article{}
	case "folder":
		result = &Folder{}
	}
	return result
}

func ToList(contentType string, obj interface{}) []contenttype.ContentTyper {
	//todo: check type first
	var result []contenttype.ContentTyper
	if contentType == "article" {
		list := obj.(*[]Article)
		for _, item := range *list {
			result = append(result, item)
		}
	} else if contentType == "folder" {
		list := obj.(*[]Folder)
		for _, item := range *list {
			result = append(result, item)
		}
	}

	return result
}
