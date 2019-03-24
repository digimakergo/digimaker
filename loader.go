package dm

import types "dm/types"

var typeLoaders = make(map[string]map[string]types.TypeLoader)

func RegisterTypeLoaders() {

	//register default loaders
	typeLoaders["contenttype"] = make(map[string]types.TypeLoader)
	typeLoaders["datatype"] = make(map[string]types.TypeLoader)

	//Get loaders
	defaultLoader := types.TypeLoaderDefault{}
	contentTypeList := defaultLoader.ContentTypeList()
	for _, identifier := range contentTypeList {
		typeLoaders["contenttype"][identifier] = &defaultLoader
	}
	fieldTypeList := defaultLoader.FieldTypeList()
	for _, identifer := range fieldTypeList {
		typeLoaders["datatype"][identifer] = &defaultLoader
	}

	//todo: load plugins loaders
}

func GetTypeLoaders() map[string]map[string]types.TypeLoader {
	if typeLoaders == nil {
		RegisterTypeLoaders()
	}
	return typeLoaders
}
