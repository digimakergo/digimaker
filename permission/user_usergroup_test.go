package permission

import (
	"dm/contenttype"
	"dm/fieldtype"
	"dm/handler"
	"dm/query"
	"fmt"
	"testing"
)

func TestUserPermission(m *testing.T) {
	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()
	err := LoadPolicies()
	fmt.Println(err)

	anonUser, err := handler.Querier().Fetch("user", query.Cond("location.id", 96))

	policyList, err := GetUserPermission(anonUser.GetCID())
	fmt.Println(policyList)
	fmt.Println(err)
	fmt.Println("anonaymouse user")
}
