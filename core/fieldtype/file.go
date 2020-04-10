//Author xc, Created on 2019-08-15 17:37
//{COPYRIGHTS}

package fieldtype

import (
	"github.com/xc/digimaker/core/util"
	"os"
	"strings"
)

type FileField struct {
	FieldValue
}

func (t *FileField) Scan(src interface{}) error {
	err := t.SetData(src, "file")
	return err
}

//implement FieldtypeHandler
type FileHandler struct{}

func (t FileHandler) Validate(input interface{}) (bool, string) {
	//todo: validate if the field exists or not
	return true, ""
}

func (t FileHandler) NewValueFromInput(input interface{}) interface{} {
	filename := input.(string)
	tempFolder := util.GetConfig("general", "upload_tempfolder", "dm")
	fileFolder := tempFolder + "/../uploaded"
	oldPath := tempFolder + "/" + filename
	newPath := fileFolder + "/" + filename
	//todo: create subfolder for it to save performance.
	os.Rename(oldPath, newPath)
	r := FileField{}
	r.Scan(filename)
	return r
}

func (t FileHandler) IsEmpty(input interface{}) bool {
	if strings.TrimSpace(input.(string)) == "" {
		return true
	}
	return false
}

func init() {
	RegisterHandler("file", FileHandler{})
}
