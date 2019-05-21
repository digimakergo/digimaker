package permission

import (
	"dm/contenttype"
	"dm/fieldtype"
	"testing"
)

func TestMain(m *testing.M) {
	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()

	GetPermissions(90)
}
