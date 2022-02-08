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

const defaultMaxSize = 10.0

type FileParameters struct {
	Format  string  `mapstructure:"format"`   //eg. "jpg, png" low case
	MaxSize float64 `mapstructure:"max_size"` //in MB, default 10MB.
}

type FileHandler struct {
	definition.FieldDef
	Params FileParameters
}

func (handler FileHandler) LoadInput(ctx context.Context, input interface{}, mode string) (interface{}, error) {
	filePath := fmt.Sprint(input)
	if filePath != "" {
		//check format
		formatParam := handler.Params.Format
		if formatParam != "" {
			formats := util.Split(handler.Params.Format)
			validFormat := false
			for _, format := range formats {
				lowerPath := strings.ToLower(filePath)
				if strings.HasSuffix(lowerPath, format) {
					validFormat = true
				}
			}
			if !validFormat {
				return nil, fieldtype.NewValidationError(fmt.Sprintf("Format not allowed. allowed are: %v", formatParam))
			}
		}

		absPath := config.VarFolder() + "/" + filePath
		info, err := os.Stat(absPath)
		if err != nil {
			return nil, fieldtype.NewValidationError("Can't find file of " + filePath)
		}

		maxSize := defaultMaxSize
		if handler.Params.MaxSize > 0 {
			maxSize = handler.Params.MaxSize
		}
		if info.Size() > int64(maxSize*1024*1024) {
			return nil, fieldtype.NewValidationError(fmt.Sprintf("File is should be less than %vM", handler.Params.MaxSize))
		}
	}
	return filePath, nil
}

//Image can be loaded from rest, or local api
func (handler FileHandler) BeforeStore(value interface{}, existing interface{}, mode string) (interface{}, error) {
	filePath := value.(string)

	//delete or new empty
	if filePath == "" {
		return "", nil
	}

	//no change
	if existing != nil && filePath == existing.(string) {
		return filePath, nil
	}

	//todo: support other file services or remote
	//todo: delete file if there is no version & after updating file

	//
	//new upload
	//
	//file path should be under temp
	temp := viper.GetString("general.upload_tempfolder")
	if !strings.HasPrefix(filePath, temp+"/") {
		return nil, fieldtype.NewValidationError("File needs to be under temp")
	}

	arr := strings.Split(filePath, "/")
	filename := arr[len(arr)-1]

	//create 2 level folder
	rand := util.RandomStr(3)
	secondLevel := string(rand)
	firstLevel := string(rand[0])

	newFolder := "file/" + firstLevel + "/" + secondLevel
	newFolderAbs := config.VarFolder() + "/" + newFolder
	_, err := os.Stat(newFolderAbs)
	if os.IsNotExist(err) {
		err = os.MkdirAll(newFolderAbs, 0775)
		if err != nil {
			return nil, err
		}
	}

	newPath := newFolder + "/" + filename
	newPathAbs := config.VarFolder() + "/" + newPath

	oldAbsPath := config.VarFolder() + "/" + filePath
	err = os.Rename(oldAbsPath, newPathAbs)

	if err != nil {
		errorMessage := "Can not move file to target " + filePath + ". error: " + err.Error()
		return nil, errors.New(errorMessage)
	}

	return newPath, nil
}

//After delete content, delete file, skip error if there is any wrong(eg. image not existing).
func (handler FileHandler) AfterDelete(value interface{}) error {
	path := config.VarFolder() + "/" + value.(string)
	err := os.Remove(path)
	if err != nil {
		message := fmt.Sprintf("Deleting file(path: %v) of %v error: %v. Deleting continued.", path, handler.FieldDef.Identifier, err.Error())
		log.Warning(message, "system")
	}

	return nil
}

func (handler FileHandler) DBField() string {
	return "VARCHAR (500) NOT NULL DEFAULT ''"
}

func init() {
	fieldtype.Register(
		fieldtype.Definition{
			Name:     "file",
			DataType: "string",
			NewHandler: func(def definition.FieldDef) fieldtype.Handler {
				params := FileParameters{}
				err := ConvertParameters(def.Parameters, &params)
				if err != nil {
					log.Warning("Wrong parameters on file:"+def.Identifier+". Ignore", "")
				}
				return FileHandler{FieldDef: def, Params: params}
			}})
}
