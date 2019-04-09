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
		value := []interface{}{expression.Value}
		return expression.Field + " " + expression.Operator + "?", value
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
