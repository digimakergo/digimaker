package permission

import (
	"context"
	"github.com/xc/digimaker/core/util/debug"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasAccessTo(t *testing.T) {
	context := debug.Init(context.Background())
	policyList, err := GetUserPolicies(2)
	assert.NotNil(t, policyList)
	assert.Nil(t, err)
	currentData := map[string]interface{}{"contenttype": "folder1"}
	result := HasAccessTo(2, "content/read", currentData, context)
	for _, item := range debug.GetDebugger(context).List {
		fmt.Println(item)
	}
	assert.Nil(t, err)
	assert.Equal(t, result, false)
}
