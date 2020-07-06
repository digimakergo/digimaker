//Author xc, Created on 2020-05-10 10:00
//{COPYRIGHTS}

package fieldtype

import (
	"errors"
	"os"
	"strings"

	"github.com/xc/digimaker/core/util"
)

//Image stores only the orginal image path.
//Thumbnail and image alias is generated in 'output' part(a kind of cache data).
type Image struct {
	String
}

func (i Image) Type() string {
	return "image"
}

//Image can be loaded from rest, or local api
func (i *Image) BeforeSaving() error {
	imagepath := i.String.String
	if imagepath != "" && imagepath != i.existing { //means there is a valid image change
		//todo: support other image services or remote image
		oldAbsPath := util.VarFolder() + "/" + imagepath

		//image path should be under temp
		temp := util.GetConfig("general", "upload_tempfolder")
		if !strings.HasPrefix(imagepath, temp) {
			return errors.New("Illegal image path.")
		}

		if _, err := os.Stat(oldAbsPath); err != nil {
			return errors.New("Can't find file on " + oldAbsPath)
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
				return err
			}
		}

		newPath := newFolder + "/" + filename
		newPathAbs := util.VarFolder() + "/" + newPath

		err = os.Rename(oldAbsPath, newPathAbs)
		if err != nil {
			errorMessage := "Can not move image to target " + imagepath + ". error: " + err.Error()
			return errors.New(errorMessage)
		}

		err = GenerateThumbnail(newPath)
		if err != nil {
			return err
		}

		i.String = String{String: newPath}
	}
	return nil
}

func GenerateThumbnail(imagePath string) error {
	varFolder := util.VarFolder()
	size := util.GetConfig("general", "image_thumbnail_size")
	thumbCacheFolder := varFolder + "/" + util.GetConfig("general", "image_thumbnail_folder")
	return util.ResizeImage(varFolder+"/"+imagePath, thumbCacheFolder+"/"+imagePath, size)
}

func init() {
	RegisterFieldType(
		FieldtypeDef{Type: "image", Value: "fieldtype.Image"},
		func() FieldTyper { return &Image{} })
}
