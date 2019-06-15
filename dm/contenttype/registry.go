package contenttype

import (
	"dm/dm/util"
)

//todo: use a better name. eg. ContentTypeMethod
type ContentTypeRegister struct {
	New     func() ContentTyper
	NewList func() interface{}
	ToList  func(obj interface{}) []ContentTyper
}

var contenttypeList = map[string]ContentTypeRegister{}

//Register a content type and store in global variable
func Register(contentType string, register ContentTypeRegister) {
	util.Log("Registering " + contentType)
	contenttypeList[contentType] = register
}

//Create new list.eg &[]Article{}
func NewList(contentType string) interface{} {
	return contenttypeList[contentType].NewList()
}

//Convert a *[]Article type(used for binding) to a slice of ContentTyper(use for more generic handling)
func ToList(contentType string, obj interface{}) []ContentTyper {
	return contenttypeList[contentType].ToList(obj)
}

//Create new content instance, eg. &Article{}
func NewInstance(contentType string) ContentTyper {
	return contenttypeList[contentType].New()
}
