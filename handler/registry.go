package handler

import (
	"dm/contenttype"
	"fmt"
)

type CTHandler interface {
	SetContent(content contenttype.ContentTyper, parentID ...int) error
}

var handlerRegistry map[string]CTHandler = map[string]CTHandler{}

func RegisterHandler(contentType string, handler CTHandler) {
	handlerRegistry[contentType] = handler
}

func GetHandler(contentType string) CTHandler {
	fmt.Println(contentType)
	fmt.Println(handlerRegistry)
	return handlerRegistry[contentType]
}
