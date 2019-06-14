//Author xc, Created on 2019-04-07 20:36
//{COPYRIGHTS}
package db

import (
	"strings"
)

var Operators = []string{">", ">=", "<", "<=", "=", "in", "like"} //todo: make it extendable in loading

type Expression struct {
	Field    string
	Operator string
	Value    interface{}
}

type Condition struct {
	Logic    string
	Children interface{} //can be []Condition or Expression(when it's the leaf) (eg. and( and( A, B ), C )
}

func (c Condition) Cond(fieldString string, value interface{}) Condition {
	cond := Cond(fieldString, value)
	return c.And(cond)
}

//And accept <cond>.And( <cond1>, <cond2> ),
//also <cond>.And( "id<", 2200 ) (same as <cond>.And( Cond( "id<", 2200 ) ))
func (c Condition) And(input interface{}, more ...interface{}) Condition {
	var result Condition
	switch input.(type) {
	case Condition:
		var arr []Condition
		for _, item := range more {
			arr = append(arr, item.(Condition))
		}
		result = combineExpression("and", c, input.(Condition), arr...)
	case string:
		value := more[0]
		result = c.And(Cond(input.(string), value)) //invoke myself with Condition type
	}
	return result
}

//Similar to And, Or accepts <cond>.Or( <cond1>, <cond2> ), also <cond>.Or( "id=", 2 )
func (c Condition) Or(input interface{}, more ...interface{}) Condition {
	var result Condition
	switch input.(type) {
	case Condition:
		var arr []Condition
		for _, item := range more {
			arr = append(arr, item.(Condition))
		}
		result = combineExpression("or", c, input.(Condition), arr...)
	case string:
		value := more[0]
		result = c.Or(Cond(input.(string), value)) //invoke myself with Condition type
	}
	return result
}

func combineExpression(operator string, input1 Condition, input2 Condition, more ...Condition) Condition {
	condition := Condition{}
	condition.Logic = operator
	var arr []Condition
	conditionArr := append(arr, input1, input2)
	if more != nil {
		conditionArr = append(conditionArr, more...)
	}
	condition.Children = conditionArr
	return condition
}

func Cond(fieldString string, value interface{}) Condition {
	condition := new(Condition)
	condition.Logic = ""
	fieldArr := separateFieldString(fieldString, value)
	condition.Children = Expression{Field: fieldArr[0], Operator: fieldArr[1], Value: value}
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

func separateFieldString(input string, value interface{}) [2]string {
	input = strings.TrimSpace(input)
	var result [2]string
	for _, operator := range Operators {
		if strings.HasSuffix(input, operator) {
			result[0] = strings.TrimSpace(strings.TrimSuffix(input, operator))
			result[1] = operator
			break
		}
	}
	//if operator is empty, it can be = if value is string/int, in if value is array
	if result[1] == "" {
		switch value.(type) {
		case string, int:
			result[0] = input
			result[1] = "="
		case []string, []int:
			result[0] = input
			result[1] = "IN"
		}
	}

	return result
}
