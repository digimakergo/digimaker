package fieldtype

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	s := &String{String: "hello"}
	err := s.Scan(nil)
	assert.Nil(t, err)
	assert.Equal(t, "", s.String)
}
