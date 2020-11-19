//Author xc, Created on 2020-05-01 16:50
//{COPYRIGHTS}
package contenttype

//User interface which is also a content type
type User interface {
	ContentTyper

	ID() int
	Username() string
	Email() string
}
