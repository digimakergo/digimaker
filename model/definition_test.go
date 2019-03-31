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

}
