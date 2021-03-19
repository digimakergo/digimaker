package tokenmanager

import (
	"context"
	"time"

	"github.com/digimakergo/digimaker/core/auth"
	"github.com/digimakergo/digimaker/core/db"
)

const tableName = "dm_token_state"

//DBTokenManager stores guid and expiry time to db and maintains them.
//todo: wirte clean up script to remove expired entries.
type DBTokenManager struct {
}

func (d DBTokenManager) Store(ctx context.Context, id string, expiry int64, claims map[string]interface{}) error {
	dbHandler := db.DBHanlder()
	tokenState := map[string]interface{}{"guid": id, "expiry": expiry}
	_, err := dbHandler.Insert(ctx, tableName, tokenState)
	return err
}

func (d DBTokenManager) Get(id string) interface{} {
	dbHandler := db.DBHanlder()
	entity := TokenState{}
	_, err := dbHandler.GetEntity(context.Background(), &entity, tableName, db.Cond("guid", id).Cond("expiry>=", time.Now().Unix()), false)
	if err != nil || entity.GUID == "" {
		return nil
	}
	return entity
}

func (d DBTokenManager) Delete(ctx context.Context, id string) error {
	dbHandler := db.DBHanlder()
	return dbHandler.Delete(ctx, tableName, db.Cond("guid", id))
}

type TokenState struct {
	GUID   string `boil:"guid" json:"guid"`
	Expiry int    `boil:"expiry" json:"expiry"`
}

func init() {
	auth.RegisterTokenManager(DBTokenManager{})
}
