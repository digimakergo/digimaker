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
	policyList, err := permission.GetUserPermission(7)
	fmt.Println(err)
	fmt.Println("Permission")
	currentData := map[string]interface{}{"contenttype": "folder1"}
	result := HasAccessTo(policyList, "content", "read", currentData, context)
	for _, item := range debug.GetDebugger(context).List {
		fmt.Print(item)
	}
	fmt.Println(result)
}
