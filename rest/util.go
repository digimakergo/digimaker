//Author xc, Created on 2019-08-25 22:51
//{COPYRIGHTS}

package rest

import (
	"dm/core/util"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	filename, err := HandleUploadFile(r, "*")
	result := ""
	if err != nil {
		w.WriteHeader(500)
		result = err.Error()
	} else {
		result = filename
	}
	w.Write([]byte(result))
}

//Upload image, return path or error
func UploadImage(w http.ResponseWriter, r *http.Request) {
	filename, err := HandleUploadFile(r, ".gif,.jpg,.jpeg,.png")
	result := ""
	if err != nil {
		w.WriteHeader(500)
		result = err.Error()
	} else {
		result = filename
	}
	w.Write([]byte(result))
}

//Handler uploaded file, return filename & error
func HandleUploadFile(r *http.Request, filetype string) (string, error) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()

	filename := handler.Filename
	//check if file type is allowed
	filetypeArr := strings.Split(filetype, ",")
	for _, extension := range filetypeArr {
		if extension != "*" && !strings.HasSuffix(filename, extension) {
			return "", errors.New("File format not allowed.")
		}
	}

	tempFolder := util.GetConfig("General", "FileUploadFolder", "dm")
	tempFolder = "/Users/xc/go/caf-prototype/src/dm/admin/web/var/upload_temp"
	//Strip file name
	reg := regexp.MustCompile("[^-A-Za-z0-9_]]")
	filename = reg.ReplaceAllString(filename, "_") //filter out all non word characters

	//Write it to temp folder
	tempFile, err := ioutil.TempFile(tempFolder, "upload-*-"+filename)
	defer tempFile.Close()
	if err != nil {
		return "", err
	}
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	tempFile.Write(fileContent)
	pathArr := strings.Split(tempFile.Name(), "/")
	tempFilename := pathArr[len(pathArr)-1]
	return tempFilename, nil
}
