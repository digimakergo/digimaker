package permission

import (
	"dm/contenttype"
	"dm/fieldtype"
	"dm/handler"
	"fmt"
	"testing"
)

func TestUserPermission(m *testing.T) {
	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()
	err := LoadPolicies()
	fmt.Println(err)

	anonUser, err := handler.Querier().FetchByID(96)

	policyList, err := GetUserPermission(anonUser.GetCID())
	fmt.Println(policyList)
	fmt.Println(err)
	fmt.Println(policyList[0].GetPolicy())
	fmt.Println("anonaymouse user")
}
