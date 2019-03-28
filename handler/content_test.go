package handler

import (
	"dm/models/orm"
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

func TestCreate(t *testing.T) {
	// data := `{"location":{"content_type": "folder", "name": "Test folder"}}`
	var c Content
	var l orm.Location
	l.Name = "Test folder"
	l.ContentType = "folder"
	l.ContentID = 0
	c.Location = &l

	var article = orm.Article{}
	article.
		// json.Unmarshal([]byte(data), &c)
		// if c.Location.Name.IsZero() {
		// 	fmt.Printf("zero")
		// }
		// err := c.Create()
		// if err != nil {
		// 	t.Error(err)
		// }
		fmt.Printf("store data")
}
