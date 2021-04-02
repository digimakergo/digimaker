package fieldtype

import (
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/log"
)

//All the definition of fieldtypes

//Definition includes a fieldtype basic information
type Definition struct {
	Name       string                            //eg. text
	BaseType   string                            //eg. string or eg "fieldtype.CustomString"
	Package    string                            //empty if there is no additional package, otherwise it's like 'mycompany.fieldtype'. Used for generating entity's import.
	NewHandler func(definition.FieldDef) Handler //callback to create new handler for this fieldtype
}

//Fieldtyper is a implementation of a fieldtype, including main logic
type Handler interface {
	//Load from input, should return the value of BaseType, (eg. int)
	ConvertInput(input interface{}) (interface{}, error)

	//Check if the input is empty.
	IsEmpty(input interface{}) bool
}

//OutputCovnerter is implemented when fieldtype needs converting when outputting
type OutputCovnerter interface {
	ConvertOuput() interface{}
}

//BeforeSaving is implemented when fieldtype has event before saving
type BeforeSaving interface {
	BeforeSave()
}

var fieldtypeMap map[string]Definition

//Register registers a fieldtype
func Register(definition Definition) {
	name := definition.Name
	if _, ok := fieldtypeMap[name]; ok {
		log.Error("Fieldtype has been previous registered: "+name+", skipped", "system")
	}
	log.Info("Registering fieldtype: " + name)
	fieldtypeMap[definition.Name] = definition
}

//GetFieldtype return a fieldtype
func GetFieldtype(fieldtype string) Definition {
	result, ok := fieldtypeMap[fieldtype]
	if ok {
		return result
	} else {
		log.Error("Field type doesn't exist: "+fieldtype, "system")
		return Definition{}
	}
}

//GetAllFieldtype get all fieldtype
func GetAllFieldtype() []string {
	result := []string{}
	for fieldtype, _ := range fieldtypeMap {
		result = append(result, fieldtype)
	}
	return result
}
