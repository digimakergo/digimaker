package permission

import (
	"dm/contenttype"
	"dm/fieldtype"
	"fmt"
	"testing"
)

func TestMain(m *testing.T) {
	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()
	err := LoadPolicies()
	fmt.Println(err)

	policyList, err := GetPermissions(90)
	fmt.Println(policyList)
}
