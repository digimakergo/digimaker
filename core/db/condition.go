//Author xc, Created on 2019-04-07 20:36
//{COPYRIGHTS}
package db

import (
	"strings"

	"github.com/digimakergo/digimaker/core/util"
)

var Operators = []string{">", ">=", "<", "==", "<=", "!=", "=", "in", "like"} //todo: make it extendable in loading
//Note: == is for join

const logicAnd = "and"
const logicOr = "or"

//Expression is a 'leaf' condition
type Expression struct {
	Field    string
	Operator string
	Value    interface{}
}

//Condition is a self contained query condition
type Condition struct {
	Logic    string
	Children interface{} //can be []Condition or Expression(when it's the leaf) (eg. and( and( A, B ), C )
	Sortby   []string
	LimitArr []int
}

//Cond is same as And(<field>, <value>) or And( Cond( <field>, <value> ) )
func (c Condition) Cond(field string, value interface{}) Condition {
	cond := Cond(field, value)
	return c.And(cond)
}

//And accepts <cond>.And( <cond1>, <cond2> ),
//also <cond>.And( "id<", 2200 ) (same as <cond>.And( Cond( "id<", 2200 ) ))
func (c Condition) And(input interface{}, more ...interface{}) Condition {
	var result Condition
	switch input.(type) {
	case Condition:
		var arr []Condition
		for _, item := range more {
			arr = append(arr, item.(Condition))
		}
		result = combineExpression(logicAnd, c, input.(Condition), arr...)
	case string:
		value := more[0]
		result = c.And(Cond(input.(string), value)) //invoke myself with Condition type
	}
	return result
}

// Or accepts <cond>.Or( <cond1>, <cond2> ), also <cond>.Or( "id=", 2 ). Similar to And
func (c Condition) Or(input interface{}, more ...interface{}) Condition {
	var result Condition
	switch input.(type) {
	case Condition:
		var arr []Condition
		for _, item := range more {
			arr = append(arr, item.(Condition))
		}
		result = combineExpression(logicOr, c, input.(Condition), arr...)
	case string:
		value := more[0]
		result = c.Or(Cond(input.(string), value)) //invoke myself with Condition type
	}
	return result
}

func (c Condition) Sort(sortby ...string) Condition {
	c.Sortby = sortby
	return c
}

func (c Condition) Limit(offset int, number int) Condition {
	c.LimitArr = []int{offset, number}
	return c
}

//combine condition like "and", "or", etc
func combineExpression(operator string, input1 Condition, input2 Condition, more ...Condition) Condition {
	condition := Condition{}
	condition.Logic = operator
	var conditions []Condition
	conditions = append(conditions, input1, input2)
	if more != nil {
		conditions = append(conditions, more...)
	}

	//filter empty condition
	var validConditions []Condition
	for _, item := range conditions {
		if item.Children != nil {
			validConditions = append(validConditions, item)
		}
	}
	if len(validConditions) == 0 {
		condition.Logic = ""
	} else {
		condition.Children = validConditions
	}
	return condition
}

//Cond creates condition like Cond("id", 1), or Cond("id", []int{1,2}) or Cond("id>", 10)
func Cond(field string, value interface{}) Condition {
	condition := new(Condition)
	condition.Logic = ""
	fieldArr := separateFieldString(field, value)
	condition.Children = Expression{Field: fieldArr[0], Operator: fieldArr[1], Value: value}
	return *condition
}

//EmptyCond creates a empty condition without expression or value
func EmptyCond() Condition {
	return Condition{}
}

func TrueCond() Condition {
	return Condition{Logic: "", Children: Expression{Field: "", Operator: "", Value: "true"}}
}

func FalseCond() Condition {
	return Condition{Logic: "", Children: Expression{Field: "", Operator: "", Value: "false"}}
}

func separateFieldString(input string, value interface{}) [2]string {
	input = strings.TrimSpace(input)
	var result [2]string
	for _, operator := range Operators {
		suffix := operator
		if operator == "in" { //todo: support more other than in
			suffix = " " + operator
		}
		if strings.HasSuffix(input, suffix) {
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

//todo: optimize - use pointers & avoid string +
func BuildCondition(cond Condition, locationColumns ...[]string) (string, []interface{}) {
	logic := cond.Logic
	if logic == "" && cond.Children == nil {
		return "", nil
	}
	if logic == "" { //if it's a expression
		expression := cond.Children.(Expression)

		//handle join condition
		//todo: fix possible sql injection issue, with more validation on field
		if expression.Operator == "==" {
			result := expression.Field + "=" + expression.Value.(string)
			return result, nil
		}

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
			//when expression has no field and operator(eg. true/false)
			if expression.Field == "" && expression.Operator == "" {
				return expression.Value.(string), nil
			}
			value = []interface{}{expression.Value}
			operatorStr = " ?"
		}
		fieldName := expression.Field
		if len(locationColumns) > 0 {
			if !(util.Contains(locationColumns[0], fieldName) || strings.Contains(fieldName, ".")) {
				fieldName = "c." + fieldName
			}
		}
		return fieldName + " " + expression.Operator + operatorStr, value
	} else {
		//If it's a container
		childrenArr := cond.Children.([]Condition)
		var expressionList []string
		var values []interface{}
		for _, subCondition := range childrenArr {
			expressionStr, currentValues := BuildCondition(subCondition, locationColumns...)
			expressionList = append(expressionList, expressionStr)
			values = append(values, currentValues...)
		}

		listStr := strings.Join(expressionList, " "+strings.ToUpper(cond.Logic)+" ")
		if len(expressionList) > 1 {
			listStr = "(" + listStr + ")"
		}
		return listStr, values
	}
}
