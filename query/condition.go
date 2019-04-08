//Author xc, Created on 2019-04-07 20:36
//{COPYRIGHTS}
package query

import "strings"

var Operators = []string{">", ">=", "<", "<=", "=", "in"} //todo: make it extendable in loading

type Operation struct {
	Field    string
	Operator string
	Value    interface{}
}

type Condition struct {
	OperationTree []Operation
	Logic         string
}

func Cond(fieldStr string, value interface{}) *Condition {
	condition := new(Condition)
	condition.Logic = ""
	fieldArr := separateFieldStr(fieldStr)
	condition.OperationTree = []Operation{Operation{Field: fieldArr[0], Operator: fieldArr[1], Value: value}}
	return condition
}

func (c *Condition) And(input ...Condition) *Condition {
	return nil
}

func (c *Condition) Or(input ...Condition) *Condition {
	return nil
}

func And(input ...Condition) *Condition {
	return nil
}

func Or(input ...Condition) *Condition {
	return nil
}

//Parentheses
func Par(input ...Condition) *[]Condition {
	return &input
}

func separateFieldStr(input string) [2]string {
	input = strings.TrimSpace(input)
	var result [2]string
	for operator := range Operators {
		if strings.HasSuffix(input, operator) {
			result[0] = strings.TrimSpace(strings.TrimSuffix(input, operator))
			result[1] = operator
			break
		}
	}
	return result
}
