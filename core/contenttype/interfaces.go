//Author xc, Created on 2019-04-01 22:00
//{COPYRIGHTS}

package contenttype

import "database/sql"

//All the content type(eg. article, folder) will implement this interface.
type ContentTyper interface {
	// Since go's embeded struct can't really inherit well from BaseContentType(eg. ID)
	// (It refers to embeded struct instance instead of override all fields by default)
	// Interface in go is like a kind of 'abstract class' when it comes to generic with data.
	// We use property of instance to declear a general ContentType. This will integrate well with orm enitty.
	/*
		ID() int
		Published() int
		Modified() int
		RemoteID() string
	*/

	//Return all fields
	//Fields() map[string]fieldtype.Fieldtyper

	GetCID() int

	GetName() string

	SetValue(identifier string, value interface{}) error

	Store(...*sql.Tx) error

	Delete(transaction ...*sql.Tx) error

	Value(identifier string) interface{}

	ContentType() string

	IdentifierList() []string

	GetLocation() *Location

	Definition(language ...string) ContentType

	GetRelations() *ContentRelationList
}

type ContentList []ContentTyper

//For enitities.
type Entitier interface {
	TableName() string
	Values() map[string]interface{}
}
