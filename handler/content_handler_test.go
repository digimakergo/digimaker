package handler

import (
	"dm/contenttype"
	"dm/contenttype/entity"
	"dm/db"
	"dm/fieldtype"
	"dm/query"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidtion(t *testing.T) {
	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()

	// Test validation1
	handler := ContentHandler{}
	params := map[string]interface{}{"title": "ff", "body": "Hello"}
	passed, result := handler.Validate("article", params)
	assert.Equal(t, passed, true)

	// Test validation2
	params = map[string]interface{}{"title": "", "body": "Hello"}
	_, result = handler.Validate("article", params)
	assert.Equal(t, result.Fields[0].Identifier, "title")

}

func TestCreate(t *testing.T) {
	// handler := ContentHandler{}
	// params := map[string]interface{}{"title": "Test " + time.Now().Format("02.01.2006 15:04"), "body": "Hello"}
	// _, result, err := handler.Create(4, "article", params)

	// params := map[string]interface{}{"title": "Test " + time.Now().Format("02.01.2006 15:04"), "summary": "Hello"}
	// _, result, err := handler.Create(4, "folder", params)
	//
	// fmt.Println(result)
	// fmt.Println(err)
}

func TestDelete(t *testing.T) {
	handler := ContentHandler{}
	err := handler.DeleteByID(40, false)
	fmt.Println(err)
	assert.Equal(t, nil, err)
}

func TestImage(t *testing.T) {
	images := &[]entity.Image{}
	handler := db.DBHanlder()
	err := handler.GetEnity("dm_image", query.Cond("1", 1), images)
	fmt.Println("images", err)
	fmt.Println(images)
}

func TestVersion(t *testing.T) {
	// handler := ContentHandler{}
	// article, _ := querier.FetchByID(6)
	// dbHanlder, _ := db.DB()
	// tx, _ := dbHanlder.Begin()
	// _, err := handler.CreateVersion(article, 1, tx)
	// tx.Commit()
	// assert.Equal( t, nil, err )
}
