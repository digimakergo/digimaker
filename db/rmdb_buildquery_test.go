package db

import (
	. "dm/query"
	"fmt"
	"testing"
)

// func TestMain(m *testing.M) {
// 	//model.LoadDefinition()
// 	m.Run()
// }

func TestBuildQuery(t *testing.T) {
	cond := Cond("id>", 12).Or(Cond("modified>=", 1111111)).
		And(Cond("section=", "content").Or(Cond("published>=", 22222)))
	fmt.Println(BuildCondition(cond))

	cond2 := Cond("id", 2).Or("id", 4).And("modified>=", 1111111)
	fmt.Println(BuildCondition(cond2))
}
