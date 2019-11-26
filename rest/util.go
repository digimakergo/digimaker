//Author xc, Created on 2019-08-25 22:51
//{COPYRIGHTS}

package rest

import (
	"dm/core/util"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
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

	tempFolder := util.GetConfig("general", "upload_tempfolder", "dm")
	// tempFolder = "/Users/xc/go/caf-prototype/src/dm/admin/web/var/upload_temp"
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

func HtmlToPDF(w http.ResponseWriter, r *http.Request) {
	//todo: permission check
	w.Header().Set("Access-Control-Allow-Origin", "*")
	html := r.PostFormValue("html")
	name := r.PostFormValue("name")
	if html == "" || name == "" {
		HandleError(errors.New("empty data"), w)
		return
	}
	result, err := htmlToPDF(html, name)
	if err != nil {
		HandleError(err, w)
		return
	}
	w.Write([]byte("var/" + result))
}

func htmlToPDF(html string, name string) (string, error) {
	tempFolder := util.GetConfig("general", "var_folder", "dm")
	uid := util.GenerateUID()
	sourceName := "/pdf/" + name + "-" + uid + ".html"
	targetName := "/pdf/" + name + "-" + uid + ".pdf"
	source := tempFolder + sourceName
	target := tempFolder + targetName
	ioutil.WriteFile(source, []byte(html), 0777)
	output, _ := exec.Command("wkhtmltopdf", "--print-media-type", source, target).Output()
	log.Println(output)
	// if err != nil {
	// 	return "", err
	// }
	return targetName, nil
}
