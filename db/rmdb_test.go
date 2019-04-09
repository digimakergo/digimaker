package db

import (
	"dm/model"
	"dm/model/entity"
	"dm/query"
	"fmt"
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
	rmdb.GetByFields("article", query.Cond("content_id", 1), &article)
	//Update remote id of the article
	fmt.Println(article)
	article.RemoteID = "5"

	err := rmdb.Update(article)
	assert.Nil(t, err)
	var article2 entity.Article
	rmdb.GetByFields("article", query.Cond("content_id", 1), &article2)
	assert.Equal(t, article2.RemoteID, "5")

}
