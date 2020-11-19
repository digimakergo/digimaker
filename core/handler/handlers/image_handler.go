//Author xc, Created on 2019-05-10 22:48
//{COPYRIGHTS}

//Package handlers implements build-in action callbacks.
package handlers

import (
	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/handler"
)

type ImageHandler struct {
}

//When creating on server side, or import.
func (ih ImageHandler) Create(content contenttype.ContentTyper, inputs handler.InputMap, parentID int) error {
	content.SetValue("parent_id", parentID)
	return nil
}

func init() {
	handler.RegisterContentTypeHandler("image", ImageHandler{})
}
