//Author xc, Created on 2019-03-27 20:00
//{COPYRIGHTS}

//Package utils provides general utils for the project.
package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/xid"
)

//UnmarshalData Load json and unmall into variable
func UnmarshalData(filepath string, v interface{}) error {
	file, err := os.Open(filepath)
	if err != nil {
		Error("Error when loading content definition: " + err.Error())
		return err
	}
	byteValue, _ := ioutil.ReadAll(file)

	err = json.Unmarshal(byteValue, v)
	if err != nil {
		Error("Error when loading definition " + filepath + ": " + err.Error())
		return err
	}

	defer file.Close()
	return nil
}

//Generate unique id. It should be cluster safe.
func GenerateUID() string {
	guid := xid.New()
	guidStr := guid.String()
	return guidStr
}

//Convert a string array to int array
func ArrayStrToInt(strArray []string) []int {
	size := len(strArray)
	var result = make([]int, size)
	for i, str := range strArray {
		result[i], _ = strconv.Atoi(str)
	}
	return result
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

//Iternate condition rules to see if all are matching.
//If there are keys in condition rules but not in realValues, match fails.
//eg. conditions: {id: 12, type:"image"}
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
				switch realValue.(type) {
				case string:
					matchResult = conditionValue == realValue
				case []string:
					matchResult = Contains(realValue.([]string), conditionValue.(string))
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
