package pongo2

import (
	"context"
	"dm/dm/contenttype"
	"dm/dm/handler"
	"dm/dm/util"
	"dm/dm/util/debug"
	"dm/sitekit"
	"dm/sitekit/niceurl"
	"fmt"
	"strings"

	"gopkg.in/flosch/pongo2.v2"
)

func dmChildren(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	querier := handler.Querier()
	parent := in.Interface().(contenttype.ContentTyper)
	context := debug.Init(context.Background())
	children, _ := querier.Children(parent, param.String(), 2, context)
	return pongo2.AsValue(children), nil
}

func dmNiceurl(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	content := in.Interface().(contenttype.ContentTyper)
	niceurl := niceurl.GenerateUrl(content)
	return pongo2.AsValue(niceurl), nil
}

func dmTplPath(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	paramArr := util.Split(param.String())
	settings := sitekit.GetSiteSettings(paramArr[1])
	fmt.Println(param)
	fmt.Println(settings)
	path := sitekit.GetContentTemplate(in.Interface().(contenttype.ContentTyper), strings.TrimSpace(paramArr[0]), settings)
	fmt.Println(path)
	return pongo2.AsValue(path), nil
}

func init() {
	pongo2.RegisterFilter("dm_children", dmChildren)
	pongo2.RegisterFilter("dm_niceurl", dmNiceurl)
	pongo2.RegisterFilter("dm_tplpath", dmTplPath)
}
