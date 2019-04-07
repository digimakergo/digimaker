package db

import (
	"dm/model"
	"dm/model/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

//
//boil.DebugMode = true
// boil.DebugWriter = os.Stderr

func TestQuery(t *testing.T) {
	model.LoadDefinition()

	rmdb := new(RMDB)
	var article entity.Article
	rmdb.GetByID("article", 2, &article)

	assert.NotNil(t, article)

}

func TestUpdate(t *testing.T) {
	model.LoadDefinition()
	rmdb := new(RMDB)

	var article2 entity.Article
	rmdb.GetByFields("article", map[string]interface{}{"content_id": 1}, &article2)
	//Update remote id of the article
	article2.RemoteID = "1"
	rmdb.Update(article2)
	var article3 entity.Article
	rmdb.GetByFields("article", map[string]interface{}{"content_id": 1}, &article3)
	assert.Equal(t, article3.RemoteID, "1")
}
