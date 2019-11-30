//Author xc, Created on 2019-08-25 22:51
//{COPYRIGHTS}

package rest

import (
	"dm/core/handler"
	"dm/core/util"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	_ "dm/sitekit/filters"

	"github.com/gorilla/mux"
	pongo2 "gopkg.in/flosch/pongo2.v2"
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

func ExportPDF(w http.ResponseWriter, r *http.Request) {
	//todo: permission check
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	id := params["id"]

	idInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(errors.New("id not int"), w)
		return
	}
	querier := handler.Querier()
	content, err := querier.FetchByID(idInt)
	if err != nil {
		HandleError(errors.New("Content not found"), w)
		return
	}

	contenttype := content.ContentType()
	tpl := pongo2.Must(pongo2.FromFile(util.HomePath() + "/templates/pdf/" + contenttype + ".html"))
	variables := map[string]interface{}{}
	variables["content"] = content

	data, err2 := tpl.ExecuteBytes(pongo2.Context(variables))
	if err2 != nil {
		HandleError(err2, w)
		return
	}

	pdfFile, err := htmlToPDF(string(data), content.GetName()+"-"+id)
	if err != nil {
		HandleError(err, w)
		return
	}

	http.Redirect(w, r, "/var/"+pdfFile, 302)

	// variables["site"] = siteIdentifier
	//
	// variables["template"] = templatePath
	// if len(matchedData) == 0 {
	// 	variables["matched_data"] = nil
	// } else {
	// 	variables["matched_data"] = matchedData[0]
	// }
	// err := tpl.ExecuteWriter(pongo2.Context(variables), w)
	//
	// w.Write([]byte(id))
	// html := r("html")
	// name := r.PostFormValue("name")
	// if html == "" || name == "" {
	// 	HandleError(errors.New("empty data"), w)
	// 	return
	// }
	// result, err := htmlToPDF(html, name)
	// if err != nil {
	// 	HandleError(err, w)
	// 	return
	// }
	// w.Write([]byte("var/" + result))
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
	output, _ := exec.Command("wkhtmltopdf", "--javascript-delay", "1000", "-L", "0mm", "-R", "0mm", "-T", "0mm", "-B", "0mm", "--print-media-type", source, target).Output()
	log.Println(output)
	// if err != nil {
	// 	return "", err
	// }
	return targetName, nil
}

func init() {
	RegisterRoute("/util/uploadfile", UploadFile)
	RegisterRoute("/util/uploadimage", UploadImage)

	RegisterRoute("/pdf/{id}", ExportPDF)
}
