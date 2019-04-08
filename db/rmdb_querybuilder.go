//Author xc, Created on 2019-04-08 22:13
//{COPYRIGHTS}
package db

import (
	. "dm/query"
	"strings"
)

//todo: optimize - use pointers & avoid string +
func BuildCondition(cond Condition) string {
	logic := cond.Logic
	if logic == "" {
		expression := cond.Children.(Expression)
		return expression.Field + " " + expression.Operator + "?"
	} else {
		childrenArr := cond.Children.([]Condition)
		var list []string
		for _, subCondition := range childrenArr {
			expressionStr := BuildCondition(subCondition)
			list = append(list, expressionStr)
		}

		listStr := strings.Join(list, " "+cond.Logic+" ")
		str := "(" + listStr + ")"
		return str
	}
}
