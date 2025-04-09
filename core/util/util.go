//Author xc, Created on 2019-03-27 20:00
//{COPYRIGHTS}

//Package utils provides general utils for the project.
package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/digimakergo/digimaker/core/log"
	"github.com/google/uuid"
	"github.com/spf13/viper"

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

func IsIdentifier(input string) bool {
	reg, _ := regexp.Compile("^[a-z][a-z0-9_]+$")
	valid := reg.Match([]byte(input))
	return valid
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
// 4) conditions: {id:11, type: "image"} target: {id:[11, 12], type:nil} //nil will be treated as pass
func MatchCondition(conditions map[string]interface{}, target map[string]interface{}) (bool, []string) {
	matchLog := []string{}
	for key, conditionValue := range conditions {
		matchResult := false
		realValue, ok := target[key]
		if !ok {
			matchResult = false
			matchLog = append(matchLog, "Mismatch since key "+key+" doesn't exist in target.")
		} else {
			if realValue == nil {
				matchResult = true //if target has key but nil value, treat same as pass
				matchLog = append(matchLog, "Matched on "+key+"since the real value is nil")
			} else if str, ok := conditionValue.(string); ok && str == "*" {
				matchResult = true
			} else {
				switch conditionValue.(type) {
				case []interface{}:
					for _, item := range conditionValue.([]interface{}) {
						//[]string or []int
						errorLog := ""
						matchResult, errorLog = matchItem(item, realValue)
						if matchResult {
							break
						}
						if errorLog != "" {
							matchLog = append(matchLog, errorLog)
						}
					}
				default:
					if conditionValueMap, ok := conditionValue.(map[string]interface{}); ok {
						if conditionSlice, ok := conditionValueMap["subset"]; ok {
							//real value should be subset of condition value, 'subset' support []string for now.
							matchSubset := true
							for _, realValueItem := range realValue.([]string) {
								conditionSliceStr := []string{}
								for _, v := range conditionSlice.([]interface{}) {
									conditionSliceStr = append(conditionSliceStr, v.(string))
								}
								matched, matchItemLog := matchItem(realValueItem, conditionSliceStr)
								if matchItemLog != "" {
									matchLog = append(matchLog, matchItemLog)
								}
								if !matched {
									matchSubset = false
									break
								}

							}
							matchResult = matchSubset
						} else {
							matchLog = append(matchLog, "map only support subset")
						}
					} else {
						errorLog := ""
						matchResult, errorLog = matchItem(conditionValue, realValue)
						if errorLog != "" {
							matchLog = append(matchLog, errorLog)
						}
					}
				}
			}
			if !matchResult {
				matchLog = append(matchLog, "Mismatch on "+key+", expecting: "+fmt.Sprint(conditionValue)+", real: "+fmt.Sprint(realValue))
			} else {
				matchLog = append(matchLog, "Matched on "+key)
			}
		}

		if !matchResult {
			return false, matchLog
		}
	}
	return true, matchLog
}

//condition value: int|string
//target value: int|[]int|string|[]string
func matchItem(current interface{}, targetValue interface{}) (bool, string) {
	result := false
	if _, ok := current.(float64); ok {
		current = int(current.(float64))
	}

	switch current.(type) {
	case string:
		switch targetValue.(type) {
		case string:
			result = current.(string) == targetValue.(string)
		case []string:
			result = Contains(targetValue.([]string), current.(string))
		default:
			return false, "Target is not string/[]string"
		}
	case int:
		switch targetValue.(type) {
		case int:
			result = current.(int) == targetValue.(int)
		case []int: //real value contains condition int
			result = ContainsInt(targetValue.([]int), current.(int))
		default:
			return false, "Target is not int/[]int"
		}
	default:
		return false, "unsupported current type:" + fmt.Sprintf("%T", current)
	}
	return result, ""
}

//Split with triming space. "," is the default separator if no seperator is provided.
func Split(str string, seperator ...string) []string {
	sep := ""
	str = strings.TrimSpace(str)
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

var sendMailHandler func(mail MailMessage) error

type MailMessage struct {
	To          []string //multi to in one mail, for name,use "name<email>" format
	Bcc         []string //bcc
	From        string
	Subject     string
	Body        string
	Attachments []string
}

func HandleSendMail(sendMail func(mail MailMessage) error) {
	sendMailHandler = sendMail
}

// Send mail with possible attachments.
// Note: From in most case(depending on handler) will be from config file
// Security: attachment needs to be checked to make sure it can not send any file
func SendFullMail(mailMessage MailMessage) error {
	//if it's not set up, ignore sending
	if sendMailHandler == nil {
		log.Error("Mail sending function is not registered", "")
		return nil
	}
	return sendMailHandler(mailMessage)
}

//Simple sending mail without attachment
func SendMail(to []string, subject, body string, bcc ...string) error {
	mailMessage := MailMessage{To: to, Subject: subject, Body: body, Bcc: bcc}
	return SendFullMail(mailMessage)
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
	path := viper.GetString("general.vipsthumbnail")
	if path == "" {
		log.Error("vipsthumbnail not found on path: "+path, "")
		return errors.New("vipsthumbnail not found")
	}
	_, err := exec.Command(path, args...).CombinedOutput()

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

//WashPath remove .. make sure it can only go under, not upp
func SecurePath(path string) string {
	reg := regexp.MustCompile(`\.\.+`)
	return reg.ReplaceAllString(path, "")
}
