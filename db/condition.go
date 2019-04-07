//Author xc, Created on 2019-04-07 20:36
//{COPYRIGHTS}
package db //todo: make a query package since it's a layer of query

type Condition struct{}

func Cond(input ...interface{}) *Condition {
	cond := new(Condition)
	return cond
}
