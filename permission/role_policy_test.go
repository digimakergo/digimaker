package permission

import (
	"dm/contenttype"
	"dm/fieldtype"
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()
	err := LoadPolicies()
	fmt.Println(err)

	GetPermissions(90)
}
