//Author xc, Created on 2019-05-10 22:48
//{COPYRIGHTS}
package handlers

import (
	"dm/contenttype"
	"dm/handler"
)

type ImageHandler struct {
}

func (ih ImageHandler) SetContent(content contenttype.ContentTyper, parentID ...int) error {
	content.SetValue("menu_id", parentID[0]) //todo: validate more
	return nil
}

func init() {
	handler.RegisterHandler("image", ImageHandler{})
}
