package query

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	//model.LoadDefinition()
	m.Run()
}

func TestCond(t *testing.T) {
	cond := Cond("id>", 11111111)
	assert.Equal(t, (cond.Children.(Expression).Operator), ">")

	cond2 := Cond(" modified  > ", 22222222)
	cond3 := Cond("published >", 3333333)
	andCond := cond.And(cond2, cond3)
	fmt.Println(andCond)

	orCond := cond.Or(cond2, cond3).And(cond2, cond3)
	var1 := orCond.Children.([]Condition)[2]
	assert.Equal(t, var1.Children.(Expression).Field, "published")
	fmt.Println(orCond)
}

func TestContinueCond(t *testing.T) {
	cond := Cond("id>", 111111).Or(Cond("modified<", 22222)).Cond("section=", "c")
	fmt.Println(cond)
	assert.Equal(t, cond.Logic, "and")
	cond1 := Cond("id<", 11111).Or("id>", 222).And("section", "c")
	fmt.Println(cond1)
	assert.Equal(t, cond1.Children.([]Condition)[1].Children.(Expression).Field, "section")
}
