package models

import (
	"encoding/json"
	"testing"
)

func TestStore(t *testing.T) {
	jsonData := []byte(`
    {
        "Location": { "name":"Test1" },
        "Fields": {}
    }`)
	var content Content
	json.Unmarshal(jsonData, &content)
	if content.Store() != nil {
		t.Errorf("Can not publish content")
	}

}
