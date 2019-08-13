package permission

import (
	"dm/core/contenttype"
	"dm/core/fieldtype"
	"fmt"
	"testing"
)

func TestMain(m *testing.T) {
	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()

	policyList, _ := GetPermissions(90)
	fmt.Println(policyList)
}
