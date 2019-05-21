//Author xc, Created on 2019-04-08 22:13
//{COPYRIGHTS}
package db

import (
	. "dm/query"
	"strings"
)

//todo: optimize - use pointers & avoid string +
func BuildCondition(cond Condition) (string, []interface{}) {
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
		return expression.Field + " " + expression.Operator + operatorStr, value
	} else {
		childrenArr := cond.Children.([]Condition)
		var list []string
		var values []interface{}
		for _, subCondition := range childrenArr {
			expressionStr, currentValues := BuildCondition(subCondition)
			list = append(list, expressionStr)
			values = append(values, currentValues...)
		}

		listStr := strings.Join(list, " "+cond.Logic+" ")
		str := "(" + listStr + ")"
		return str, values
	}
}
