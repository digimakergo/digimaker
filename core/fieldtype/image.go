//Author xc, Created on 2020-05-10 10:00
//{COPYRIGHTS}

package fieldtype

import (
	"errors"
	"os"
	"strings"

	"github.com/xc/digimaker/core/log"
	"github.com/xc/digimaker/core/util"
)

type Image struct {
	String
	changed bool
}

func (i Image) Type() string {
	return "image"
}

func (i *Image) LoadFromInput(input interface{}) error {
	original := i.String.String
	str := i.String
	err := str.LoadFromInput(input)
	if err == nil {
		i.String = str
		if original != i.String.String {
			i.changed = true
		}
	}
	return err
}

//Image can be loaded from rest, or local api
func (i *Image) BeforeSaving() error {
	filepath := i.String.String
	log.Error("GOOD", "")
	if i.changed && filepath != "" {
		//todo: support other image services or remote image
		oldAbsPath := util.VarFolder() + "/" + filepath
		if _, err := os.Stat(oldAbsPath); err != nil {
			return errors.New("Can't find file on " + oldAbsPath)
		}
		arr := strings.Split(filepath, "/")
		filename := arr[len(arr)-1]

		//create 2 level folder
		rand := util.RandomStr(3)
		secondLevel := string(rand)
		firstLevel := string(rand[0])

		newFolder := "images/" + firstLevel + "/" + secondLevel
		newFolderAbs := util.VarFolder() + "/" + newFolder
		_, err := os.Stat(newFolderAbs)
		if os.IsNotExist(err) {
			err = os.MkdirAll(newFolderAbs, 0775)
			if err != nil {
				return err
			}
		}
		newPath := newFolder + "/" + filename
		newPathAbs := util.VarFolder() + "/" + newPath

		//todo: create thumbnail
		err = os.Rename(oldAbsPath, newPathAbs)
		if err != nil {
			errorMessage := "Can not move image to target " + filepath + ". error: " + err.Error()
			return errors.New(errorMessage)
		}
		i.String = String{String: newPath}
		log.Error(i.String.String, "")
	}
	return nil
}

func init() {
	RegisterFieldType(
		FieldtypeDef{Type: "image", Value: "fieldtype.Image"},
		func() FieldTyper { return &Image{} })
}
