package permission

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserPolicies(t *testing.T) {
	policyList, err := GetUserPolicies(context.Background(), 2)
	assert.NotNil(t, policyList)
	assert.Nil(t, err)
}

func ExampleHasAccessTo() {
	currentData := map[string]interface{}{"contenttype": "folder"}

	//2 is a member
	result := HasAccessTo(context.Background(), 2, "content/read", currentData)
	fmt.Println(result)
	//output: true
}
