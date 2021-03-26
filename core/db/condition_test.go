package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	//contenttype.LoadDefinition()
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

	falseCond := FalseCond()
	fmt.Println(BuildCondition(falseCond))
}

func ExampleCond() {
	//cond
	cond := Cond("id>", 10)

	output, _ := BuildCondition(cond)
	fmt.Println(output)
	//output: id > ?
}

func ExampleCond_mulitiCond() {
	cond := Cond("modified >", 22222)

	//Here Cond() is the same as And()
	andCond := cond.Cond("published >", 3333333)

	output, _ := BuildCondition(andCond)
	fmt.Println(output)
	//output: (modified > ? AND published > ?)
}

func ExampleCond_andCondition() {
	cond := Cond(" modified  > ", 22222222)
	cond2 := Cond("published >", 3333333)

	//and
	andCond := cond.And(cond2)

	output, _ := BuildCondition(andCond)
	fmt.Println(output)
	//output: (modified > ? AND published > ?)
}

func ExampleCond_andExpression() {
	cond := Cond("modified >", 22222)

	//and
	andCond := cond.And("published >", 3333333) //Here is the same as cond.And( db.Cond( "published >", 3333333 ) )

	output, _ := BuildCondition(andCond)
	fmt.Println(output)
	//output: (modified > ? AND published > ?)
}

func ExampleCond_or() {
	cond := Cond(" modified  > ", 22222222)
	cond2 := Cond("published >", 3333333)

	//or
	orCond := cond.Or(cond2)

	output, _ := BuildCondition(orCond)
	fmt.Println(output)

	//output: (modified > ? OR published > ?)
}

func ExampleCond_orAnd() {
	cond := Cond(" modified  > ", 22222222)
	cond2 := Cond("published >", 3333333)

	//or then and
	orCond := cond.Or(cond2).And("id>", 1)

	output, _ := BuildCondition(orCond)
	fmt.Println(output)

	//output: ((modified > ? OR published > ?) AND id > ?)
}

func ExampleEmptyCond() {
	cond := EmptyCond().Cond("author", "1")
	built, _ := BuildCondition(cond)
	fmt.Println(built)
	//Output: author = ?
}

func TestContinueCond(t *testing.T) {
	cond := Cond("id>", 111111).Or(Cond("modified<", 22222)).Cond("section=", "c")
	fmt.Println(cond)
	assert.Equal(t, cond.Logic, "and")
	cond1 := Cond("id<", 11111).Or("id>", 222).And("section", "c")
	fmt.Println(cond1)
	assert.Equal(t, cond1.Children.([]Condition)[1].Children.(Expression).Field, "section")
}

func TestBuildQuery(t *testing.T) {
	cond := Cond("id>", 12).Or(Cond("modified>=", 1111111)).
		And(Cond("section=", "content").Or(Cond("published>=", 22222)))
	fmt.Println(BuildCondition(cond))

	cond2 := Cond("id", 2).Or("id", 4).And("modified>=", 1111111)
	fmt.Println(BuildCondition(cond2))

	cond3 := Cond("id", []string{"1", "2"})
	fmt.Println(BuildCondition(cond3))
}
