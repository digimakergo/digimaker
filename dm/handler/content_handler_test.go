package handler

import (
	"context"
	"dm/admin/entity"
	"dm/dm"
	"dm/dm/contenttype"
	"dm/dm/db"
	"dm/dm/util/debug"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var ctx context.Context

func TestMain(m *testing.M) {
	dm.Bootstrap("/Users/xc/go/caf-prototype/src/dm/test")
	fmt.Println("Test starting..")
	ctx = debug.Init(context.Background())
	m.Run()
}

func TestValidtion(t *testing.T) {

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

var contentCreated contenttype.ContentTyper

func TestCreate(t *testing.T) {
	ctx := debug.Init(context.Background())
	handler := ContentHandler{Context: ctx}
	// // params := map[string]interface{}{"title": "Test " + time.Now().Format("02.01.2006 15:04"), "body": "Hello"}
	// // _, result, err := handler.Create(4, "article", params)
	//
	params := map[string]interface{}{"title": "Test " + time.Now().Format("02.01.2006 15:04"), "summary": "Hello"}
	result, validation, err := handler.Create("folder", params, 1)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, true, validation.Passed())

	contentCreated, _ = Querier().FetchByContentID("folder", result.GetCID())
}

func TestUpdate1(t *testing.T) {
	handler := ContentHandler{Context: ctx}
	folder, _ := querier.FetchByID(contentCreated.GetLocation().ID)
	inputs := map[string]interface{}{"summary": "updated"}
	pass, _, err := handler.Update(folder, inputs)
	assert.Nil(t, err)
	assert.Equal(t, true, pass)
}

func TestCreateImage(t *testing.T) {
	ctx := debug.Init(context.Background())
	handler := ContentHandler{Context: ctx}

	params := map[string]interface{}{"title": "Test " + time.Now().Format("02.01.2006 15:04"), "path": "Hello"}
	result, validation, err := handler.Create("image", params, 1)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, true, validation.Passed())
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

func TestDelete(t *testing.T) {
	handler := ContentHandler{}
	handler.Context = ctx
	id := contentCreated.GetLocation().ID
	fmt.Println(id)
	err := handler.DeleteByID(id, false)
	assert.Equal(t, nil, err)
}

func TestQueryImage(t *testing.T) {
	images := &[]entity.Image{}
	handler := db.DBHanlder()
	err := handler.GetEntity("dm_image", db.Cond("1", 1), images)
	assert.Nil(t, err)
	assert.NotNil(t, images)
}
