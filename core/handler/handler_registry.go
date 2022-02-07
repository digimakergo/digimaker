//Author xc, Created on 2019-05-10 22:20
//{COPYRIGHTS}

//Package handler implements content related actions(eg.create/edit/delete/import) and callback mechanism while handling.
package handler

import (
	"strconv"
	"strings"

	"github.com/digimakergo/digimaker/core/config"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/util"
)

//Content type handler registeration list
var contentTypeHandlerList map[string]interface{} = map[string]interface{}{}

func RegisterContentTypeHandler(contentType string, handler interface{}) {
	handlerType := []string{}
	if _, ok := handler.(ContentTypeHandlerValidate); ok {
		handlerType = append(handlerType, "validate")
	}

	if _, ok := handler.(ContentTypeHandlerCreate); ok {
		handlerType = append(handlerType, "create")
	}
	if _, ok := handler.(ContentTypeHandlerUpdate); ok {
		handlerType = append(handlerType, "update")
	}
	if _, ok := handler.(ContentTypeHandlerDelete); ok {
		handlerType = append(handlerType, "delete")
	}

	log.Info("Registering contenttype handler for " + contentType + ", implemented: " + strings.Join(handlerType, ","))
	contentTypeHandlerList[contentType] = handler
}

func GetContentTypeHandler(contentType string) interface{} {
	return contentTypeHandlerList[contentType]
}

//Operation handler registeration list
var operationHandlerList []OperationHandler = []OperationHandler{}

func RegisterOperationHandler(handler OperationHandler) {
	log.Info("Registering operation handler " + handler.Identifier)
	operationHandlerList = append(operationHandlerList, handler)
}

//Get operation handler list based on rules defined in operation_handler.json/yaml
//target is the real vaule the condition matches to
//target should not include 'event' key since it's in the parameter already.
func GetOperationHandlerByCondition(event string, target map[string]interface{}) ([]OperationHandler, []string) {
	//todo: preserve order in the config so matched event will be called from top to down
	viperHandler := config.GetViper("handler")
	handlers := viperHandler.GetStringMap("handlers")
	//todo: test this after config change
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
