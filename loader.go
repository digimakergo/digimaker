//Author xc, Created on 2019-03-25 20:00
//{COPYRIGHTS}

package dm

import "dm/fieldtype"

//TypeLoader is a interface for plugin to register content type, data type, etc.
//Basically one typerloader is one plugin.
type TypeLoader interface {
	Instance(extendedType string, identifier string) interface{}
	FieldTypeList() []string
	ContentTypeList() []string
}

var typeLoaders = make(map[string]map[string]TypeLoader)

func RegisterTypeLoaders() {

	//register default loaders
	typeLoaders["contenttype"] = make(map[string]TypeLoader)
	typeLoaders["datatype"] = make(map[string]TypeLoader)

	//Get loaders
	defaultLoader := fieldtype.TypeLoaderDefault{}
	contentTypeList := defaultLoader.ContentTypeList()
	for _, identifier := range contentTypeList {
		typeLoaders["contenttype"][identifier] = &defaultLoader
	}
	fieldTypeList := defaultLoader.FieldTypeList()
	for _, identifer := range fieldTypeList {
		typeLoaders["datatype"][identifer] = &defaultLoader
	}

	//todo: load plugins loaders via plugin mechanism, in additional to TyperLoaderDefault
}

func GetTypeLoaders() map[string]map[string]TypeLoader {
	if typeLoaders == nil {
		RegisterTypeLoaders()
	}
	return typeLoaders
}
