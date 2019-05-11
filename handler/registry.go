//Author xc, Created on 2019-05-10 22:20
//{COPYRIGHTS}

package handler

import (
	"dm/util"
	"strconv"
	"strings"
)

//Content type handler registeration list
var contentTypeHandlerList map[string]ContentTypeHandler = map[string]ContentTypeHandler{}

func RegisterContentTypeHandler(contentType string, handler ContentTypeHandler) {
	contentTypeHandlerList[contentType] = handler
}

func GetContentTypeHandler(contentType string) ContentTypeHandler {
	return contentTypeHandlerList[contentType]
}

//Operation handler registeration list
var operationHandlerList []OperationHandler = []OperationHandler{}

func RegisterOperationHandler(handler OperationHandler) {
	operationHandlerList = append(operationHandlerList, handler)
}

//Get operation handler list based on rules defined in operation_handler.json/yaml
//target is the real vaule the condition matches to
//target should not include 'event' key since it's in the parameter already.
func GetOperationHandlerByCondition(event string, target map[string]interface{}) ([]OperationHandler, []string) {
	//todo: preserve order in the config so matched event will be called from top to down
	handlers := util.GetConfigSectionI("handlers", "operation_handler")
	target["event"] = event
	matchLog := []string{}
	keys := []string{}
	for key, _ := range target {
		keys = append(keys, key)
	}
	matchLog = append(matchLog, "Matching with target with keys: "+strings.Join(keys, ", "))
	result := []OperationHandler{}
	for identifier, conditions := range handlers {
		conditionMap := conditions.(map[string]interface{})
		var (
			matchResult     bool
			currentMatchLog []string
		)

		if _, ok := conditionMap["event"]; !ok {
			matchResult = false
			currentMatchLog = []string{"No 'event' defined in the condition."}
		} else {
			matchResult, currentMatchLog = util.MatchCondition(conditionMap, target)
		}
		matchLog = append(matchLog, "Match "+identifier+": "+
			strconv.FormatBool(matchResult)+
			". Matching detail: "+strings.Join(currentMatchLog, ", "))
		if matchResult {
			for _, handler := range operationHandlerList {
				if handler.Identifier == identifier {
					result = append(result, handler)
				}
			}
		}
	}
	return result, matchLog
}
