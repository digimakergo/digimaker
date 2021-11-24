//Author xc, Created on 2019-08-25 22:51
//{COPYRIGHTS}

package rest

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/digimakergo/digimaker/core/permission"
	"github.com/digimakergo/digimaker/core/util"

	"github.com/gorilla/mux"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	section := util.GetConfigSectionI("rest")
	needAuth := section["upload_file_auth"].(bool)
	if needAuth {
		userId := CheckUserID(r.Context(), w)
		if userId == 0 || !permission.HasAccessTo(r.Context(), userId, "util/upload") {
			HandleError(errors.New("Need authorization"), w)
			return
		}
	}

	filename, err := HandleUploadFile(r, "*")
	result := ""
	if err != nil {
		HandleError(err, w)
		return
	}

	result = filename
	WriteResponse(result, w)
}

//Upload image, return path or error
func UploadImage(w http.ResponseWriter, r *http.Request) {
	userId := CheckUserID(r.Context(), w)
	if userId == 0 {
		return
	}
	filename, err := HandleUploadFile(r, ".gif,.jpg,.jpeg,.png")
	result := ""
	if err != nil {
		HandleError(err, w)
		return
	}

	result = filename
	WriteResponse(result, w)
}

//Handler uploaded file, return filename & error
func HandleUploadFile(r *http.Request, filetype string) (string, error) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()

	filename := strings.ToLower(handler.Filename)
	//check if file type is allowed
	fileAllowed := false
	filetypeArr := strings.Split(filetype, ",")
	for _, extension := range filetypeArr {
		if extension == "*" || strings.HasSuffix(filename, extension) {
			fileAllowed = true
			break
		}
	}
	if !fileAllowed {
		return "", errors.New("File format not allowed.")
	}

	tempFolder := util.GetConfig("general", "upload_tempfolder")
	tempFolderAbs := util.VarFolder() + "/" + tempFolder

	//Strip file name
	reg := regexp.MustCompile("[^-A-Za-z0-9_.]")
	filename = reg.ReplaceAllString(filename, "_") //filter out all non word characters

	//Write it to temp folder
	tempFile, err := ioutil.TempFile(tempFolderAbs, "upload-*-"+filename)
	defer tempFile.Close()
	if err != nil {
		return "", err
	}
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	tempFile.Chmod(0664)
	tempFile.Write(fileContent)
	pathArr := strings.Split(tempFile.Name(), "/")
	tempFilename := pathArr[len(pathArr)-1]
	return tempFolder + "/" + tempFilename, nil
}

func GetAllowedLimitations(w http.ResponseWriter, r *http.Request) {
	userId := CheckUserID(r.Context(), w)
	if userId == 0 {
		return
	}

	params := mux.Vars(r)
	operation := params["operation"]
	operation = strings.ReplaceAll(operation, "_", "/")

	allowedOperations := util.GetConfigArr("permission", "rest_allowed_operations", "dm")
	if !util.Contains(allowedOperations, operation) {
		HandleError(errors.New("Operation not allowed"), w, 403)
		return
	}

	_, limits, err := permission.GetUserAccess(r.Context(), userId, operation)
	if err != nil {
		HandleError(err, w)
		return
	}
	WriteResponse(limits, w)
}

func init() {
	RegisterRoute("/util/uploadfile", UploadFile, "POST")
	RegisterRoute("/util/uploadimage", UploadImage, "POST")
	RegisterRoute("/util/limitations/{operation}", GetAllowedLimitations)
}
