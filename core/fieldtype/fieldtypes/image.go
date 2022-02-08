package fieldtypes

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/digimakergo/digimaker/core/config"
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/spf13/viper"
)

type ImageParameter struct {
	Format string `mapstructure:"format"` //eg. "jpg, png" low case
}

type ImageHandler struct {
	definition.FieldDef
}

func (handler ImageHandler) LoadInput(ctx context.Context, input interface{}, mode string) (interface{}, error) {
	//todo: check image format
	str := fmt.Sprint(input)
	return str, nil
}

//Image can be loaded from rest, or local api
func (handler ImageHandler) BeforeStore(ctx context.Context, value interface{}, existing interface{}, mode string) (interface{}, error) {
	imagePath := value.(string)

	if imagePath == "" {
		return "", nil //todo: remove existing if there is
	}

	if existing == nil || existing != nil && imagePath != existing.(string) { //means there is a valid image change
		//todo: support other image services or remote image

		params := ImageParameter{}
		err := ConvertParameters(handler.FieldDef.Parameters, &params)
		if err != nil {
			return existing, err
		}

		//check format
		if params.Format != "" {
			formats := util.Split(params.Format)
			validFormat := false
			for _, format := range formats {
				lowImagePath := strings.ToLower(imagePath)
				if strings.HasSuffix(lowImagePath, format) {
					validFormat = true
				}
			}
			if !validFormat {
				return existing, fmt.Errorf("Format not allowed. allowed are: %v", params.Format)
			}
		}

		oldAbsPath := config.VarFolder() + "/" + imagePath

		//image path should be under temp
		temp := viper.GetString("general.upload_tempfolder")

		if _, err := os.Stat(oldAbsPath); err != nil {
			return nil, errors.New("Can't find file on " + oldAbsPath)
		}

		arr := strings.Split(imagePath, "/")
		filename := arr[len(arr)-1]

		//create 2 level folder
		rand := util.RandomStr(3)
		secondLevel := string(rand)
		firstLevel := string(rand[0])

		newFolder := "images/" + firstLevel + "/" + secondLevel
		newFolderAbs := config.VarFolder() + "/" + newFolder
		_, err = os.Stat(newFolderAbs)
		if os.IsNotExist(err) {
			err = os.MkdirAll(newFolderAbs, 0775)
			if err != nil {
				return nil, err
			}
		}

		newPath := newFolder + "/" + filename
		newPathAbs := config.VarFolder() + "/" + newPath

		underTemp := strings.HasPrefix(imagePath, temp)
		if underTemp {
			err = os.Rename(oldAbsPath, newPathAbs)
		} else {
			err = os.Link(oldAbsPath, newPathAbs) //todo: use better copy
		}
		if err != nil {
			errorMessage := "Can not copy/move image to target " + imagePath + ". error: " + err.Error()
			return nil, errors.New(errorMessage)
		}

		err = GenerateThumbnail(newPath)
		if err != nil {
			return nil, err
		}
		return newPath, nil
	} else {
		if existing != nil {
			return existing, nil
		} else {
			return "", nil
		}
	}
}

//After delete content, delete image&thumbnail, skip error if there is any wrong(eg. image not existing).
func (handler ImageHandler) AfterDelete(ctx context.Context, value interface{}) error {
	path := config.VarFolder() + "/" + value.(string)
	err := os.Remove(path)
	if err != nil {
		message := fmt.Sprintf("Deleting image(path: %v) of %v error: %v. Deleting continued.", path, handler.FieldDef.Identifier, err.Error())
		log.Warning(message, "system")
	}

	thumbnail := ThumbnailFolder() + "/" + path
	err = os.Remove(thumbnail)
	if err != nil {
		message := fmt.Sprintf("Deleting image thumnail(path: %v) of %v error: %v. Deleting continued.", path, handler.FieldDef.Identifier, err.Error(), "system")
		log.Warning(message, "system")
	}
	return nil
}

func (handler ImageHandler) DBField() string {
	return "VARCHAR (500) NOT NULL DEFAULT ''"
}

func GenerateThumbnail(imagePath string) error {
	varFolder := config.VarFolder()
	size := viper.GetString("general.image_thumbnail_size")
	thumbCacheFolder := ThumbnailFolder()
	return util.ResizeImage(varFolder+"/"+imagePath, thumbCacheFolder+"/"+imagePath, size)
}

func ThumbnailFolder() string {
	thumbFolder := config.VarFolder() + "/" + viper.GetString("general.image_thumbnail_folder")
	return thumbFolder
}

func init() {
	fieldtype.Register(
		fieldtype.Definition{
			Name:     "image",
			DataType: "string",
			NewHandler: func(def definition.FieldDef) fieldtype.Handler {
				return ImageHandler{FieldDef: def}
			}})
}
