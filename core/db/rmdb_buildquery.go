//Author xc, Created on 2019-04-08 22:13
//{COPYRIGHTS}
package db

import (
	"strings"

	"github.com/xc/digimaker/core/util"
)

//todo: optimize - use pointers & avoid string +
func BuildCondition(cond Condition, locationColumns ...[]string) (string, []interface{}) {
	logic := cond.Logic
	if logic == "" {
		expression := cond.Children.(Expression)
		value := []interface{}{}
		operatorStr := ""
		switch expression.Value.(type) {
		//when value is string slice
		case []string:
			for _, item := range expression.Value.([]string) {
				value = append(value, item)
			}
			operatorArr := []string{}
			for _ = range value {
				operatorArr = append(operatorArr, "?")
			}
			operatorStr = " (" + strings.Join(operatorArr, ",") + ")"
			//when value is int slice
		case []int:
			for _, item := range expression.Value.([]int) {
				value = append(value, item)
			}
			operatorArr := []string{}
			for _ = range value {
				operatorArr = append(operatorArr, "?")
			}
			operatorStr = " (" + strings.Join(operatorArr, ",") + ")"
			//when value is string/int
		default:
			value = []interface{}{expression.Value}
			operatorStr = " ?"
		}
		fieldName := expression.Field
		if len(locationColumns) > 0 && fieldName != "1" {
			if !(util.Contains(locationColumns[0], fieldName) || strings.Contains(fieldName, ".")) {
				fieldName = "content." + expression.Field
			}
		}
		return fieldName + " " + expression.Operator + operatorStr, value
	} else {
		childrenArr := cond.Children.([]Condition)
		var expressionList []string
		var values []interface{}
		for _, subCondition := range childrenArr {
			expressionStr, currentValues := BuildCondition(subCondition, locationColumns...)
			expressionList = append(expressionList, expressionStr)
			values = append(values, currentValues...)
		}

		listStr := strings.Join(expressionList, " "+cond.Logic+" ")
		str := "(" + listStr + ")"
		return str, values
	}
}
