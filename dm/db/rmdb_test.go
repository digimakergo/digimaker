//Author xc, Created on 2019-04-27 20:33
//{COPYRIGHTS}
package db

import (
	"dm/dm/query"
	"fmt"
	"testing"
)

func TestInsert(t *testing.T) {
	// rmdb := RMDB{}
	//
	// //test without transaction.
	// values := map[string]interface{}{
	// 	"description": "Test1"}
	// result, err := rmdb.Insert("dm_relation", values)
	// fmt.Println(result)
	// assert.Equal(t, nil, err)
	//
	// //Test with transation
	// values = map[string]interface{}{
	// 	"description": "Test2"}
	//
	// database, err := DB()
	// assert.Equal(t, nil, err)
	// tx, err := database.Begin()
	// assert.Equal(t, nil, err)
	// _, err = rmdb.Insert("dm_relation", values, tx)
	// tx.Commit()
	// if err != nil {
	// 	tx.Rollback()
	// }
}

func TestCount(t *testing.T) {
	rmdb := RMDB{}
	count, err := rmdb.Count("dm_location", query.Cond("1", 1))
	fmt.Println(count, err)
}
