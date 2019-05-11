//Author xc, Created on 2019-05-10 22:20
//{COPYRIGHTS}

package handler

import (
	"dm/util"
	"fmt"
	"strconv"
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

func GetOperationHandlerList() []OperationHandler {
	return operationHandlerList
}

//Get operation handler list based on rules defined in operation_handler.json/yaml
//target is the real vaule the condition matches to
func GetOperationHandlerByCondition(target map[string]interface{}) ([]OperationHandler, []string) {
	handlers := util.GetConfigSectionI("handlers", "operation_handler")
	matchLog := []string{}
	result := []OperationHandler{}
	for identifier, conditions := range handlers {
		matchLog = append(matchLog, "matching "+identifier)
		//match
		matchResult, currentMatchLog := util.MatchCondition(conditions.(map[string]interface{}), target)
		matchLog = append(matchLog, "matching result: "+
			strconv.FormatBool(matchResult)+
			"matching detail:"+fmt.Sprint(currentMatchLog))
		if matchResult {
			for _, handler := range operationHandlerList {
				if handler.Identifer() == identifier {
					result = append(result, handler)
				}
			}
		}
	}
	return result, matchLog
}
