package tokenmanager

import (
	"fmt"

	"github.com/xc/digimaker/core/db"
	"github.com/xc/digimaker/rest"
)

const tableName = "dm_token_state"

type DBTokenManager struct {
}

func (d DBTokenManager) Store(id string, expiry int64) error {
	dbHandler := db.DBHanlder()
	tokenState := map[string]interface{}{"guid": id, "expiry": expiry}
	sid, err := dbHandler.Insert(tableName, tokenState)
	fmt.Println(sid)
	return err
}

func (d DBTokenManager) Get(id string) interface{} {
	dbHandler := db.DBHanlder()
	entity := TokenState{}
	err := dbHandler.GetEntity(tableName, db.Cond("guid", id), nil, &entity)
	if err != nil || entity.GUID == "" {
		return nil
	}
	return entity
}

func (d DBTokenManager) Delete(id string) error {
	return nil
}

type TokenState struct {
	GUID   string `boil:"guid" json:"guid"`
	Expiry int    `boil:"expiry" json:"expiry"`
}

func init() {
	rest.RegisterTokenManager(DBTokenManager{})
}
