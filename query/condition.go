//Author xc, Created on 2019-04-07 20:36
//{COPYRIGHTS}
package query

import "strings"

var Operators = []string{">", ">=", "<", "<=", "=", "in"} //todo: make it extendable in loading

type Expression struct {
	Field    string
	Operator string
	Value    interface{}
}

type Condition struct {
	Logic    string
	Children interface{} //can be []Condition or Expression(when it's the leaf) (eg. and( and( A, B ), C )
}

func Cond(fieldStr string, value interface{}) Condition {
	condition := new(Condition)
	condition.Logic = ""
	fieldArr := separateFieldStr(fieldStr)
	condition.Children = Expression{Field: fieldArr[0], Operator: fieldArr[1], Value: value}
	return *condition
}

func (c Condition) And(input Condition, more ...Condition) Condition {
	return combineExpression("and", c, input, more...)
}

func (c Condition) Or(input Condition, more ...Condition) Condition {
	return combineExpression("or", c, input, more...)
}

func combineExpression(operator string, input1 Condition, input2 Condition, more ...Condition) Condition {
	condition := new(Condition)
	condition.Logic = operator
	var arr []Condition
	conditionArr := append(arr, input1, input2)
	if more != nil {
		conditionArr = append(conditionArr, more...)
	}
	condition.Children = conditionArr
	return *condition
}

/*
func And(input ...Condition) *Condition {
	return nil
}

func Or(input ...Condition) *Condition {
	return nil
}
*/

//Parentheses
func Par(input ...Condition) *[]Condition {
	return &input
}

func separateFieldStr(input string) [2]string {
	input = strings.TrimSpace(input)
	var result [2]string
	for _, operator := range Operators {
		if strings.HasSuffix(input, operator) {
			result[0] = strings.TrimSpace(strings.TrimSuffix(input, operator))
			result[1] = operator
			break
		}
	}
	return result
}
