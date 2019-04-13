package entity

import "dm/model"

func NewContentType(contentType string) model.ContentTyper {
	var result model.ContentTyper
	switch contentType {
	case "article":
		result = Article{}
	case "folder":
		result = Folder{}
	}
	return result
}
