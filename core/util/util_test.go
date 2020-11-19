package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtil(t *testing.T) {
	str := "Hello world  test: dd"
	assert.Equal(t, "hello-world-test-dd", NameToIdentifier(str))
}

func ExampleNameToIdentifier() {
	result := GetStrVar("{this} is {good}")
	fmt.Println(result[0])
	//Output: this
}

func ExampleReplaceStrVar() {
	result := ReplaceStrVar("{this} is {value}", map[string]string{"this": "test", "value": "test1"})
	fmt.Println(result)
	//Output: test is test1
}

func ExampleUpperName() {
	result := UpperName("hello_world")
	fmt.Println(result)
	//Output: HelloWorld
}
