package handler

import (
	"context"
	"dm/permission"
	"dm/util/debug"
	"fmt"
	"testing"
)

func TestHasAccessTo(t *testing.T) {
	context := debug.Init(context.Background())
	policyList, err := permission.GetUserPolicies(7)
	fmt.Println(policyList)
	fmt.Println(err)
	fmt.Println("Permission")
	currentData := map[string]interface{}{"contenttype": "folder1"}
	result, _ := HasAccessTo(7, "content", "read", currentData, context)
	for _, item := range debug.GetDebugger(context).List {
		fmt.Print(item)
	}
	fmt.Println(result)
}
