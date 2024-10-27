//Author xc, Created on 2019-08-25 22:51
//{COPYRIGHTS}

package rest

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/digimakergo/digimaker/core/config"
	"github.com/digimakergo/digimaker/core/permission"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/spf13/viper"

	"github.com/gorilla/mux"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	needAuth := viper.GetBool("rest.upload_file_auth")
	if needAuth {
		userId := CheckUserID(r.Context(), w)
		if userId == 0 || !permission.HasAccessTo(r.Context(), userId, "util/upload") {
			HandleError(errors.New("Need authorization"), w)
			return
		}
	}

	service := r.URL.Query().Get("service")
	result := ""

	if service == "image" {
		imageResult, err := HandleUploadImage(r)
		if err != nil {
			HandleError(err, w)
			return
		}
		result = imageResult
	} else {
		filename, err := HandleUploadFile(r, "*")
		if err != nil {
			HandleError(err, w)
			return
		}

		result = filename
	}
	WriteResponse(result, w)
}

// Upload image, return path or error
func UploadImage(w http.ResponseWriter, r *http.Request) {
	userId := CheckUserID(r.Context(), w)
	if userId == 0 {
		HandleError(errors.New("Need login"), w)
		return
	}

	result, err := HandleUploadImage(r)
	if err != nil {
		HandleError(err, w)
		return
	}

	WriteResponse(result, w)
}

func HandleUploadImage(r *http.Request) (string, error) {
	filename, err := HandleUploadFile(r, ".gif,.jpg,.jpeg,.png,.pdf")
	result := ""
	if err != nil {
		return "", err
	}
	if strings.HasSuffix(filename, ".pdf") {
		newImage, _ := strings.CutSuffix(filename, ".pdf")
		newImage = newImage + ".jpg"
		pdfPath := config.PathWithVar(filename)
		err := util.ResizeImage(pdfPath, config.PathWithVar(newImage), "1000x>")
		if err != nil {
			return "", err
		}
		os.Remove(pdfPath)
		result = newImage
	} else {
		result = filename
	}
	return result, nil
}

// Handler uploaded file, return filename & error
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

	tempFolder := viper.GetString("general.upload_tempfolder")
	tempFolderAbs := config.PathWithVar(tempFolder)

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

	allowedOperations := viper.GetStringSlice("permission.rest_allowed_operations")
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
