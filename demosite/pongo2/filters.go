package pongo2

import (
	"context"
	"dm/dm/contenttype"
	"dm/dm/handler"
	"dm/dm/util/debug"
	"dm/niceurl"

	"github.com/flosch/pongo2"
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

func init() {
	pongo2.RegisterFilter("dm_children", dmChildren)
	pongo2.RegisterFilter("dm_niceurl", dmNiceurl)

}
