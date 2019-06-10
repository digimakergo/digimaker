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

	policyList, _ := GetPermissions(90)
	fmt.Println(policyList)
}
