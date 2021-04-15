//Author xc, Created on 2019-04-01 22:00
//{COPYRIGHTS}

package contenttype

import (
	"context"
	"database/sql"

	"github.com/digimakergo/digimaker/core/definition"
)

//All the content type(eg. article, folder) will implement this interface.
type ContentTyper interface {
	GetCID() int

	GetID() int

	GetName() string

	SetValue(identifier string, value interface{}) error

	Store(ctx context.Context, transaction ...*sql.Tx) error

	Delete(ctx context.Context, transaction ...*sql.Tx) error

	Value(identifier string) interface{}

	ContentType() string

	IdentifierList() []string

	GetLocation() *Location

	Definition(language ...string) definition.ContentType

	ToMap() map[string]interface{}
}

type GetRelations interface {
	GetRelations() ContentRelationList
}

type ContentList []ContentTyper
