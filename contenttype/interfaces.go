//Author xc, Created on 2019-04-01 22:00
//{COPYRIGHTS}

package contenttype

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
	//Fields() map[string]fieldtype.Fielder

	Values() map[string]interface{}

	//Visit  field dynamically
	//Field(name string) interface{}

	//Visit all attribute dynamically including Fields + internal attribute eg. id, parent_id.
	//Attr(name string) interface{}
}

type ContentList []ContentTyper

//For enitities.
type Entitier interface {
	TableName() string
	Values() map[string]interface{}
}
