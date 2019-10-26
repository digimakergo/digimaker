package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtil(t *testing.T) {
	str := "Hello world  test: dd"
	assert.Equal(t, "hello-world-test-dd", NameToIdentifier(str))
}

func TestNameToIdentifier(t *testing.T) {
	result := GetStrVar("{this} is {good}")
	assert.Equal(t, result[0], "this")
}

func TestReplaceStrVar(t *testing.T) {
	result := ReplaceStrVar("{this} is {value}", map[string]string{"this": "test", "value": "test1"})
	assert.Equal(t, "test is test1", result)
}
