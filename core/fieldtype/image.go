package fieldtype

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/util"
)

type ImageHandler struct {
	definition.FieldDef
}

func (handler ImageHandler) LoadInput(input interface{}, mode string) (interface{}, error) {
	//todo: check image format
	str := fmt.Sprint(input)
	return str, nil
}

//Image can be loaded from rest, or local api
func (i ImageHandler) BeforeSaving(value interface{}, existing interface{}, mode string) (interface{}, error) {
	imagepath := value.(string)

	if imagepath != "" && imagepath != existing.(string) { //means there is a valid image change
		//todo: support other image services or remote image
		oldAbsPath := util.VarFolder() + "/" + imagepath

		//image path should be under temp
		temp := util.GetConfig("general", "upload_tempfolder")

		if _, err := os.Stat(oldAbsPath); err != nil {
			return nil, errors.New("Can't find file on " + oldAbsPath)
		}

		arr := strings.Split(imagepath, "/")
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
				return nil, err
			}
		}

		newPath := newFolder + "/" + filename
		newPathAbs := util.VarFolder() + "/" + newPath

		underTemp := strings.HasPrefix(imagepath, temp)
		if underTemp {
			err = os.Rename(oldAbsPath, newPathAbs)
		} else {
			err = os.Link(oldAbsPath, newPathAbs)
		}
		if err != nil {
			errorMessage := "Can not copy/move image to target " + imagepath + ". error: " + err.Error()
			return nil, errors.New(errorMessage)
		}

		err = GenerateThumbnail(newPath)
		if err != nil {
			return nil, err
		}
		return newPath, nil
	} else {
		//todo: remove the existing thumbnails
		return "", nil
	}

}

func GenerateThumbnail(imagePath string) error {
	varFolder := util.VarFolder()
	size := util.GetConfig("general", "image_thumbnail_size")
	thumbCacheFolder := varFolder + "/" + util.GetConfig("general", "image_thumbnail_folder")
	return util.ResizeImage(varFolder+"/"+imagePath, thumbCacheFolder+"/"+imagePath, size)
}

func init() {
	Register(
		Definition{
			Name:     "text",
			DataType: "string",
			NewHandler: func(def definition.FieldDef) Handler {
				return ImageHandler{FieldDef: def}
			}})
}
