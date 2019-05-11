//Author xc, Created on 2019-03-27 20:00
//{COPYRIGHTS}

//Package utils provides general utils for the project
package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
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
		Error("Error when loading datatype definition: " + err.Error())
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

//Convert like "hello_world" to "HelloWorld"
func UpperName(input string) string {
	arr := strings.Split(input, "_")
	for i := range arr {
		arr[i] = strings.Title(arr[i])
	}
	return strings.Join(arr, "")
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
	matchResult := true
	matchingLog := []string{}
	for key, conditionValue := range conditions {
		realValue, ok := target[key]
		if ok {
			switch conditionValue.(type) {
			case int, string:
				matchResult = matchResult && conditionValue == conditionValue
			case []int, []string:
				matchResult = matchResult && Contains(conditionValue.([]string), realValue.(string))
			}
			if !matchResult {
				matchingLog = append(matchingLog, "mismatch on "+key+",expecting: "+fmt.Sprint(conditionValue)+", real: "+fmt.Sprint(realValue))
			}
		} else {
			matchResult = false
			matchingLog = append(matchingLog, "mismatch since "+key+"doesn't exist in target.")
		}
		if !matchResult {
			break
		}
	}
	return matchResult, matchingLog
}
