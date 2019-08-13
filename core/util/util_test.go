package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtil(t *testing.T) {
	str := "Hello world  test: dd"
	assert.Equal(t, "hello-world-test-dd", NameToIdentifier(str))
}
