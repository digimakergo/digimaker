//Author xc, Created on 2019-08-15 17:37
//{COPYRIGHTS}

package fieldtype

import (
	"dm/core/util"
	"os"
	"strings"
)

type ImageField struct {
	FieldValue
}

func (t *ImageField) Scan(src interface{}) error {
	err := t.SetData(src, "image")
	return err
}

//implement FieldtypeHandler
type ImageHandler struct{}

func (t ImageHandler) Validate(input interface{}) (bool, string) {
	//todo: validate if the field exists or not
	return true, ""
}

func (t ImageHandler) NewValueFromInput(input interface{}) interface{} {
	filename := input.(string)
	//todo: implement differnet way to handle different files(eg. upload, upload to cluster, dropbox, image service, etc)
	tempFolder := util.GetConfig("general", "upload_tempfolder", "dm")
	imageFolder := tempFolder + "/../uploaded"
	oldPath := tempFolder + "/" + filename
	newPath := imageFolder + "/" + filename
	//todo: create subfolder for it to save performance.
	//todo: create thumbnail
	os.Rename(oldPath, newPath)
	r := ImageField{}
	r.Scan(filename)
	return r
}

func (t ImageHandler) IsEmpty(input interface{}) bool {
	if strings.TrimSpace(input.(string)) == "" {
		return true
	}
	return false
}

func init() {
	RegisterHandler("image", ImageHandler{})
}
