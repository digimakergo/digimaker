package db

import (
	"dm/model"
	"dm/model/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	model.LoadDefinition()
	m.Run()
}

//
//boil.DebugMode = true
// boil.DebugWriter = os.Stderr

func TestQuery(t *testing.T) {

	rmdb := new(RMDB)
	var article entity.Article
	rmdb.GetByID("article", 2, &article)

	assert.NotNil(t, article)

}

func TestUpdate(t *testing.T) {
	rmdb := new(RMDB)

	var article entity.Article
	rmdb.GetByFields("article", map[string]interface{}{"content_id": 1}, &article)
	//Update remote id of the article
	article.RemoteID = "4"
	//rmdb.Update(article)
	err := rmdb.Update(article)
	assert.Nil(t, err)
	var article2 entity.Article
	rmdb.GetByFields("article", map[string]interface{}{"content_id": 1}, &article2)
	assert.Equal(t, article2.RemoteID, "4")
}
