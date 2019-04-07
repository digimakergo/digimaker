//Author xc, Created on 2019-04-01 22:00
//{COPYRIGHTS}
package db

import (
	"context"
	"dm/model"
	"dm/model/entity"
	"dm/util"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
)

// Implement DBEntitier
type RMDB struct{}

//Query by ID
func (rmdb *RMDB) GetByID(contentType string, id int, content model.ContentTyper) error {
	return rmdb.GetByFields(contentType, map[string]interface{}{"id": id}, content)
}

//Query to fill in contentTyper. Use reference in content parameter.
//It fill in with nil if nothing found(no error returned in this case)
//  var content model.Article
//  rmdb.GetByFields("article", map[string]interface{}{"id": 12}, content)
//
func (*RMDB) GetByFields(contentType string, fields interface{}, content model.ContentTyper) error {
	db, err := DB()
	if err != nil {
		return errors.New("Error when connecting db: " + err.Error())
	}

	contentTypeDef := model.ContentTypeDefinition[contentType]
	tableName := contentTypeDef.TableName

	//get condition string for fields
	fieldStr := ""
	var values []interface{}
	for name, value := range fields.(map[string]interface{}) {
		_, isLocationField := model.LocationFields[name]
		nameWithTable := "c." + name
		if isLocationField {
			nameWithTable = "l." + name
		}
		fieldStr += "AND " + nameWithTable + "=?"
		values = append(values, value)
	}

	sql := `SELECT * FROM dm_location l, ` + tableName + ` c
		 		 	WHERE l.content_id=c.id
								AND l.content_type= '` + contentType + `'
							  ` + fieldStr

	util.Debug("db", sql)
	err = queries.Raw(sql, values...).Bind(context.Background(), db, content)

	if err != nil {
		return errors.New("Error when query table: " + tableName + " " + err.Error())
	}
	return nil
}

func (*RMDB) Update(article entity.Article) error {
	db, _ := DB()
	_, err := article.Update(context.Background(), db, boil.Infer())
	if err != nil {
		return nil
	}
	return nil
}

func (*RMDB) GetEntities() {

}
