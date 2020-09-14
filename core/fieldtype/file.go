package fieldtype

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/xc/digimaker/core/util"
)

type File struct {
	String
}

func (f File) Type() string {
	return "file"
}

func (f *File) LoadFromInput(input interface{}, params FieldParameters) error {
	err := f.String.LoadFromInput(input, params)
	if err != nil {
		return err
	}

	path := f.String.String
	if path == "" {
		return nil
	}

	index := strings.LastIndex(path, ".")
	if index == -1 {
		return errors.New("Empty file format")
	}
	fileSuffix := path[index+1:]
	fmt.Println(fileSuffix)

	if format, ok := params["format"]; ok {
		formatList := format.([]interface{})
		formatListS := []string{}
		for _, f := range formatList {
			formatListS = append(formatListS, f.(string))
		}

		if !util.Contains(formatListS, fileSuffix) {
			return errors.New("Invalid file format: " + fileSuffix)
		}
	}
	return nil
}

func (f *File) BeforeSaving() error {
	tempPath := f.String.String
	if tempPath != "" && tempPath != f.existing {
		oldAbsPath := util.VarFolder() + "/" + tempPath

		//image path should be under temp
		temp := util.GetConfig("general", "upload_tempfolder")
		if !strings.HasPrefix(tempPath, temp) {
			return errors.New("Illegal image path.")
		}

		if _, err := os.Stat(oldAbsPath); err != nil {
			return errors.New("Can't find file on " + oldAbsPath)
		}

		arr := strings.Split(tempPath, "/")
		filename := arr[len(arr)-1]

		//create 2 level folder
		rand := util.RandomStr(3)
		secondLevel := string(rand)
		firstLevel := string(rand[0])

		newFolder := "documents/" + firstLevel + "/" + secondLevel
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
			errorMessage := "Can not move file to target " + tempPath + ". error: " + err.Error()
			return errors.New(errorMessage)
		}

		f.String = String{String: newPath}
	}
	return nil
}

func init() {
	RegisterFieldType(
		FieldtypeDef{Type: "file", Value: "fieldtype.File"},
		func() FieldTyper { return &File{} })
}
