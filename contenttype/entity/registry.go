package entity

import "dm/contenttype"

//todo: use a better name. eg. ContentTypeMethod
type ContentTypeRegister struct {
	New     func() contenttype.ContentTyper
	NewList func() interface{}
	ToList  func(obj interface{}) []contenttype.ContentTyper
}

var contenttypeList = map[string]ContentTypeRegister{}

//Register a content type and store in global variable
func Register(contentType string, register ContentTypeRegister) {
	contenttypeList[contentType] = register
}

//Create new list.eg &[]Article{}
func NewList(contentType string) interface{} {
	return contenttypeList[contentType].NewList()
}

func ToList(contentType string, obj interface{}) []contenttype.ContentTyper {
	return contenttypeList[contentType].ToList(obj)
}

//Create new content instance, eg. &Article{}
func NewInstance(contentType string) contenttype.ContentTyper {
	return contenttypeList[contentType].New()
}
