package pongo2

import (
	"context"
	"dm/core/contenttype"
	"dm/core/handler"
	"dm/core/util"
	"dm/core/util/debug"
	"dm/sitekit"
	"dm/sitekit/niceurl"
	"fmt"
	"os"
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

func dmTplMatched(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	paramArr := util.Split(param.String())
	settings := sitekit.GetSiteSettings(paramArr[1])
	fmt.Println(param)
	fmt.Println(settings)
	path := sitekit.GetContentTemplate(in.Interface().(contenttype.ContentTyper), strings.TrimSpace(paramArr[0]), settings)
	fmt.Println(path)
	return pongo2.AsValue(path), nil
}

func dmTplPath(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	path := ""
	if !param.IsNil() {
		//get path under gopath
		packageName := param.String()
		path = os.Getenv("GOPATH") + "/src/" + packageName + "/templates/" + in.String()
	} else {
		//get path under current package
		path = util.HomePath() + "/templates/" + in.String()
	}

	return pongo2.AsValue(path), nil
}

func init() {
	pongo2.RegisterFilter("dm_children", dmChildren)
	pongo2.RegisterFilter("dm_niceurl", dmNiceurl)
	pongo2.RegisterFilter("dm_tpl_matched", dmTplMatched)
	pongo2.RegisterFilter("dm_tpl_path", dmTplPath)
}
