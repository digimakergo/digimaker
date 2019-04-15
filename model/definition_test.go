//Author xc, Created on 2019-03-28 21:00
//{COPYRIGHTS}

package model

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	err := LoadDefinition()
	if err != nil {
		t.Fail()
	}

	t.Log(fmt.Printf(ContentTypeDefinition["folder"].TableName + "\n"))
	t.Log(fmt.Printf(ContentTypeDefinition["article"].TableName + "\n"))

	s := "hello world"
	//var text TextDatatype = "GOod"
	fmt.Println(TextDatatype(s))
}
