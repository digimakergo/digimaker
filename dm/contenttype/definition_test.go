//Author xc, Created on 2019-03-28 21:00
//{COPYRIGHTS}

package contenttype

import (
	"dm/util"
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	//todo: remove it and in test
	var path = "/Users/xc/go/caf-prototype/src/dm/configs"
	util.SetConfigPath(path)
	err := LoadDefinition()
	if err != nil {
		t.Fail()
	}

	t.Log(fmt.Printf(contentTypeDefinition["folder"].TableName + "\n"))
	t.Log(fmt.Printf(contentTypeDefinition["article"].TableName + "\n"))

}
