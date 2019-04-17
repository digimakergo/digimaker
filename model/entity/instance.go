package entity

func NewInstance(contentType string) interface{} {
	// var result model.ContentTyper
	// switch contentType {
	// case "article":
	// 	result = Article{}
	// case "folder":
	// 	result = Folder{}
	// }
	return Article{}
}

func NewListInstance(contentType string, instance *interface{}) {
	var result interface{}
	switch contentType {
	case "article":
		result = []Article{}
	case "folder":
		result = []Folder{}
	}
	instance = &result
}
