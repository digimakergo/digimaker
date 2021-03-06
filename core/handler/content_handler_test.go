package handler

import (
	"context"
	"testing"
	"time"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/digimakergo/digimaker/test"

	"github.com/stretchr/testify/assert"
)

var ctx context.Context

func TestMain(m *testing.M) {
	ctx = test.Start()
	m.Run()
}

func TestValidate(t *testing.T) {
	// Test validation1
	handler := ContentHandler{}
	params := map[string]interface{}{"title": "ff", "body": "Hello"}
	def, _ := definition.GetDefinition("article")
	passed, result := handler.Validate("article", def.FieldMap, params, true)
	assert.Equal(t, true, passed)

	// Test validation2
	params = map[string]interface{}{"title": "", "body": "Hello"}
	_, result = handler.Validate("article", def.FieldMap, params, true)
	assert.Equal(t, "1", result.Fields["title"])

	params = map[string]interface{}{"title": nil, "body": "Hello"}
	_, result = handler.Validate("article", def.FieldMap, params, true)
	assert.Equal(t, "1", result.Fields["title"])

	params = map[string]interface{}{"body": "Hello"}
	_, result = handler.Validate("article", def.FieldMap, params, true)
	assert.Equal(t, "1", result.Fields["title"])

}

var contentCreated contenttype.ContentTyper

func TestCreate(t *testing.T) {
	handler := ContentHandler{Context: ctx}

	content, validation, err := handler.Create("article", map[string]interface{}{"title": "title only"}, 1, 3)

	assert.Nil(t, err)
	assert.Equal(t, true, validation.Passed())
	assert.Equal(t, "title only", content.Value("title").(*fieldtype.Text).FieldValue().(string))

	params := map[string]interface{}{"title": "Test " + time.Now().Format("02.01.2006 15:04"), "summary": "Hello"}
	result, validation, err := handler.Create("folder", params, 1, 3)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, true, validation.Passed())

	contentCreated, _ = Querier().FetchByContentID("folder", result.GetCID())
}

func TestUpdate(t *testing.T) {
	handler := ContentHandler{Context: ctx}
	folder, _ := querier.FetchByID(contentCreated.GetLocation().ID)

	//update will fail because of lacking of title
	inputs := map[string]interface{}{"summary": "updated"}
	pass, vResult, err := handler.Update(folder, inputs, 1)
	assert.Equal(t, false, pass)
	assert.Equal(t, "1", vResult.Fields["title"])

	//update will succeed
	inputs = map[string]interface{}{"title": "test", "summary": "updated"}
	pass, _, err = handler.Update(folder, inputs, 1)
	assert.Equal(t, true, pass)
	assert.Equal(t, nil, err)

}

func TestCreateImage(t *testing.T) {
	handler := ContentHandler{Context: ctx}

	params := map[string]interface{}{"title": "Test " + time.Now().Format("02.01.2006 15:04"), "path": "Hello"}
	result, validation, err := handler.Create("image", params, 1, 3)

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

	//Create and delete
	content, _, _ := handler.Create("folder", map[string]interface{}{"title": "delete"}, 1, 3)

	id := content.GetLocation().ID
	err := handler.DeleteByID(id, 1, false)
	assert.Equal(t, nil, err)
}
