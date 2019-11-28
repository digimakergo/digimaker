package pongo2

import (
	"context"
	"dm/core/contenttype"
	"dm/core/handler"
	"dm/core/util"
	"dm/core/util/debug"
	"dm/sitekit"
	"dm/sitekit/niceurl"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/flosch/pongo2.v2"
)

func dmChildren(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	querier := handler.Querier()
	parent := in.Interface().(contenttype.ContentTyper)
	context := debug.Init(context.Background())
	children, _, _ := querier.Children(parent, param.String(), 2, []int{}, []string{}, false, context)
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

func dmFormatTime(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	timestamp := in.Integer()
	result := time.Unix(int64(timestamp), 0).Format("02.01.2006 15:04")
	return pongo2.AsValue(result), nil
}

func dmConfig(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	filename := in.String()
	if param.IsNil() {
		err := pongo2.Error{ErrorMsg: "Need section param."}
		return pongo2.AsValue(""), &err
	}
	section := param.String()
	result := util.GetConfigSectionAll(section, filename)
	return pongo2.AsValue(result), nil
}

func dmValue(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	input := in.Interface().(map[string]interface{})
	key := param.String()
	result := input[key]
	return pongo2.AsValue(result), nil
}

func dmJson(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	input := in.String()
	var result interface{}
	err := json.Unmarshal([]byte(input), &result)
	if err != nil {
		result = nil
	}
	return pongo2.AsValue(result), nil
}

func dmStr(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	input := in.Interface()
	result := fmt.Sprintln(input)
	return pongo2.AsValue(result), nil
}

func dmNow(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	result := time.Now()
	return pongo2.AsValue(result), nil
}

func init() {
	pongo2.RegisterFilter("dm_children", dmChildren)
	pongo2.RegisterFilter("dm_niceurl", dmNiceurl)
	pongo2.RegisterFilter("dm_tpl_matched", dmTplMatched)
	pongo2.RegisterFilter("dm_tpl_path", dmTplPath)
	pongo2.RegisterFilter("dm_format_time", dmFormatTime)
	pongo2.RegisterFilter("dm_config", dmConfig)
	pongo2.RegisterFilter("dm_json", dmJson)
	pongo2.RegisterFilter("dm_str", dmStr)
	pongo2.RegisterFilter("dm_value", dmValue)
	pongo2.RegisterFilter("dm_now", dmNow)

}
