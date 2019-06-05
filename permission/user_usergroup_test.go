package permission

import (
	"dm/contenttype"
	"dm/fieldtype"
	"fmt"
	"testing"
)

func TestUserPermission(m *testing.T) {
	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()

	policyList, err := GetUserPermission(7)
	fmt.Println(policyList)
	fmt.Println(err)
	fmt.Println(policyList[0].GetPolicy())
	fmt.Println("anonaymouse user")
}
