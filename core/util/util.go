//Author xc, Created on 2019-03-27 20:00
//{COPYRIGHTS}

//Package utils provides general utils for the project.
package util

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xc/digimaker/core/log"

	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

//UnmarshalData Load json and unmall into variable
func UnmarshalData(filepath string, v interface{}) error {
	file, err := os.Open(filepath)
	if err != nil {
		log.Error("Error when loading content definition: "+err.Error(), "")
		return err
	}
	byteValue, _ := ioutil.ReadAll(file)

	err = json.Unmarshal(byteValue, v)
	if err != nil {
		log.Error("Error when loading definition "+filepath+": "+err.Error(), "")
		return err
	}

	defer file.Close()
	return nil
}

//Generate unique id with order. It should be cluster safe.
func GenerateUID() string {
	guid := xid.New()
	guidStr := guid.String()
	return guidStr
}

//Generate a guid which is completely random without order
func GenerateGUID() string {
	uuid := uuid.New()
	return uuid.String()
}

//Convert a string array to int array
func ArrayStrToInt(strArray []string) ([]int, error) {
	size := len(strArray)
	var result = make([]int, size)
	for i, str := range strArray {
		value, err := strconv.Atoi(str)
		result[i] = value
		if err != nil {
			return []int{}, err
		}
	}
	return result, nil
}

var varBrakets = []string{"{", "}"}

//Get variable from defined brakets. eg "{name} is {realname}" will get ["name", "realname"]
func GetStrVar(str string) []string {
	left := varBrakets[0]
	right := varBrakets[1]
	r, _ := regexp.Compile(left + `(\w+)` + right)
	match := r.FindAllStringSubmatch(str, -1)
	result := []string{}
	for i := range match {
		result = append(result, match[i][1])
	}
	return result
}

//Replace variable with values in string
func ReplaceStrVar(str string, values map[string]string) string {
	result := str
	for key := range values {
		value := values[key]
		old := varBrakets[0] + key + varBrakets[1]
		result = strings.ReplaceAll(result, old, value)
	}
	return result
}

//Convert like "hello_world" to "HelloWorld"
func UpperName(input string) string {
	arr := strings.Split(input, "_")
	for i := range arr {
		arr[i] = strings.Title(arr[i])
	}
	return strings.Join(arr, "")
}

func InterfaceToStringArray(input []interface{}) []string {
	result := make([]string, len(input))
	for i, value := range input {
		result[i] = value.(string)
	}
	return result
}

//convert name lie "Hello world.?" to "hello-world"
func NameToIdentifier(input string) string {
	lowerStr := strings.ToLower(strings.TrimSpace(input))
	reg, _ := regexp.Compile("[^a-z0-9]+")
	result := reg.ReplaceAllString(lowerStr, "-")
	return result
}

//Strip unregular string to avoid sql injection.
//Note this is not used to strip whole sql, but phrase of sql.(eg. ORDER BY ...), not applicable for values.
func StripSQLPhrase(str string) string {
	reg, _ := regexp.Compile("[^a-z0-9A-Z., _]+")
	result := reg.ReplaceAllString(str, "")
	return result
}

//Iternate condition rules to see if all are matching.
//If there are keys in condition rules but not in realValues, match fails. * mean always all
//eg.1) conditions: {id: 12, type:"image"} or {id:[11,12], type:["image", "article"]} target: {id:12,type:"article"}
// 2) conditions: {id: [11,12], type:"image" } target: {id:[12, 13], type: ["image", "article"]}
// 3) conditions: {id:11, type: "*"} target: {id:[11, 12], type:"image"}
func MatchCondition(conditions map[string]interface{}, target map[string]interface{}) (bool, []string) {
	matchResult := false
	matchLog := []string{}
	for key, conditionValue := range conditions {
		realValue, ok := target[key]
		if ok {
			switch conditionValue.(type) {
			case int:
				switch realValue.(type) {
				case int:
					matchResult = conditionValue == realValue
				case []int: //real value contains condition int
					matchResult = ContainsInt(realValue.([]int), conditionValue.(int))
				}
			case string:
				if conditionValue.(string) == "*" {
					matchResult = true
				} else {
					switch realValue.(type) {
					case string:
						matchResult = conditionValue == realValue
					case []string:
						matchResult = Contains(realValue.([]string), conditionValue.(string))
					}
				}
			case []interface{}:
				for _, item := range conditionValue.([]interface{}) {
					if _, ok := item.(string); ok {
						if _, ok := realValue.(string); ok {
							matchResult = item.(string) == realValue.(string)
						} else {
							matchLog = append(matchLog, "Target value should be string.")
						}
					}
					if _, ok := item.(int); ok {
						if _, ok := realValue.(int); ok {
							matchResult = item.(int) == realValue.(int)
						} else {
							matchLog = append(matchLog, "Target value should be int")
						}
					}
					if matchResult {
						break
					}
				}
			}

			if !matchResult {
				matchLog = append(matchLog, "Mismatch on "+key+", expecting: "+fmt.Sprint(conditionValue)+", real: "+fmt.Sprint(realValue))
			} else {
				matchLog = append(matchLog, "Matched on "+key)
			}
		} else {
			matchResult = false
			matchLog = append(matchLog, "Mismatch since key "+key+" doesn't exist in target.")
		}
		if !matchResult {
			break
		}
	}
	return matchResult, matchLog
}

//Split with triming space. "," is the default separator if no seperator is provided.
func Split(str string, seperator ...string) []string {
	sep := ""
	if len(seperator) == 0 {
		sep = ","
	} else {
		sep = seperator[0]
	}
	arr := strings.Split(str, sep)
	for i, value := range arr {
		arr[i] = strings.TrimSpace(value)
	}
	return arr
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14) //todo: make cost configable
	return string(hash), err
}

func MatchPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

//todo: rewrite this.
func SendMail(to []string, subject, body string) error {
	from := GetConfig("general", "send_from", "dm")
	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")

	addr := "127.0.0.1:25"
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Mail(r.Replace(from)); err != nil {
		return err
	}
	for i := range to {
		to[i] = r.Replace(to[i])
		if err = c.Rcpt(to[i]); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	msg := "To: " + strings.Join(to, ",") + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

//RandomStr generate a random string. no number only small letters.
func RandomStr(n int) []byte {
	rand.Seed(time.Now().UTC().UnixNano())
	characters := "abcdefghijklmnopqrstuvwxyz"
	str := []byte("")
	for i := 0; i < n; i++ {
		j := rand.Intn(len(characters))
		str = append(str, characters[j])
	}
	return str
}

//Generate a resized image
func ResizeImage(from string, to string, size string) error {
	folder := filepath.Dir(to)
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.MkdirAll(folder, 0775)
		if err != nil {
			return err
		}
	}

	var args = []string{
		"--size", size,
		"--output", to,
		from,
	}
	path, err := exec.LookPath("vipsthumbnail")
	if err != nil {
		return err
	}
	cmd := exec.Command(path, args...)
	err = cmd.Run()

	if err != nil {
		log.Error("Can not resize image. args:"+fmt.Sprint(args), "")
		return errors.New("Can not resize image " + from + ": " + err.Error())
	}
	return nil
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func init() {
	path, err := exec.LookPath("vipsthumbnail")
	if err != nil {
		log.Warning("vipsthumbnail not found for image proceeding.", "")
	} else {
		log.Info("Vipsthumbnail found in " + path)
	}
}
