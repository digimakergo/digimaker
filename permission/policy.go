//Author xc, Created on 2019-05-22 22:10
//{COPYRIGHTS}

package permission

type Permission struct {
	Module     string
	Action     []string
	Limitation interface{}
}

type Policy struct {
	Identifier  string
	Name        string
	LimitedTo   []string
	Permissions []Permission
}

func GetUserPermissions(userID int) {
	//
}
