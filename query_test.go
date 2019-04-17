package dm

import (
	"dm/db"
	"dm/def"
	"dm/handler"
	"dm/model/entity"
	. "dm/query"
	"dm/util"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	def.LoadDefinition()
	m.Run()
}

//
//boil.DebugMode = true
// boil.DebugWriter = os.Stderr

func TestQuery(t *testing.T) {

	rmdb := new(db.RMDB)
	var article entity.Article
	rmdb.GetByID("article", 2, &article)

	assert.NotNil(t, article)

}

func TestUpdate(t *testing.T) {
	rmdb := new(db.RMDB)

	var article entity.Article
	rmdb.GetByFields("article", Cond("content_id", 1), &article)
	//Update remote id of the article
	fmt.Println(article)
	uid := util.GenerateUID()
	println(uid)
	article.RemoteID = uid
	err := article.Store()
	fmt.Println(err)

	/*
		id, error := rmdb.Insert("dm_article", map[string]interface{}{"modified": 231213})
		if error != nil {
			fmt.Println(id, error.Error())
		}
	*/

	err = rmdb.Update(article.TableName(), map[string]interface{}{"body": "test" + time.Now().String()}, Cond("id", 1))
	assert.Nil(t, err)
	var article2 entity.Article
	rmdb.GetByFields("article", Cond("content_id", 1), &article2)

	assert.Equal(t, article2.RemoteID, uid)

	//insert
	// article3 := new(entity.Article)
	// article3.RemoteID = "5555555"
	// article3.Store()

	var articles []entity.Article
	handler.Query.List("article", Cond("1", "1"), &articles)
	fmt.Println(articles)

}
