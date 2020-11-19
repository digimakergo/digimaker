package pongo2

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/handler"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/digimakergo/digimaker/sitekit"
	"github.com/digimakergo/digimaker/sitekit/niceurl"
	"golang.org/x/text/message"

	"gopkg.in/flosch/pongo2.v2"
)

func dmChildren(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	querier := handler.Querier()
	parent := in.Interface().(contenttype.ContentTyper)
	children, _, _ := querier.Children(parent, param.String(), 1, db.EmptyCond(), []int{}, []string{}, false, context.Background())
	return pongo2.AsValue(children), nil
}

func dmWashNumber(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	language := param.String()
	p := message.NewPrinter(message.MatchLanguage(language))
	result := p.Sprintln(in.Interface())
	return pongo2.AsValue(result), nil
}

func dmFetch(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	querier := handler.Querier()
	id := in.Interface().(int)
	content, _ := querier.FetchByID(id)
	return pongo2.AsValue(content), nil
}

func dmParent(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	querier := handler.Querier()
	content := in.Interface().(contenttype.ContentTyper)
	parentID := content.Value("parent_id").(int)
	parent, _ := querier.FetchByID(parentID)
	return pongo2.AsValue(parent), nil
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
		path = util.AbsHomePath() + "/templates/" + in.String()
	}

	return pongo2.AsValue(path), nil
}

func dmFormatTime(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	timestamp := in.Integer()
	result := time.Unix(int64(timestamp), 0).Format("02.01.2006 15:04")
	return pongo2.AsValue(result), nil
}

func dmFormatNumber(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	str := in.String()
	arr := []byte(str)
	spliter := []byte(param.String())[0]
	result := []byte{}
	count := len(str)
	div := count - count/3*3
	for i := 0; i < div; i++ {
		result = append(result, arr[i])
	}

	for j := count / 3; j >= 1; j-- {
		if len(result) > 0 {
			result = append(result, spliter, arr[count-j*3], arr[count-j*3+1], arr[count-j*3+2])
		} else {
			result = append(result, arr[count-j*3], arr[count-j*3+1], arr[count-j*3+2])
		}
	}
	resultStr := string(result)
	return pongo2.AsValue(resultStr), nil
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
	var result interface{}
	switch in.Interface().(type) {
	case []interface{}:
		input := in.Interface().([]interface{})
		key := param.Integer()
		result = input[key]
	default:
		input := in.Interface().(map[string]interface{})
		key := param.String()
		result = input[key]
	}

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

func dmSplit(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	result := strings.Split(in.String(), param.String())

	return pongo2.AsValue(result), nil
}

func init() {
	pongo2.RegisterFilter("dm_children", dmChildren)
	pongo2.RegisterFilter("dm_parent", dmParent)
	pongo2.RegisterFilter("dm_fetch", dmFetch)
	pongo2.RegisterFilter("dm_niceurl", dmNiceurl)
	pongo2.RegisterFilter("dm_tpl_matched", dmTplMatched)
	pongo2.RegisterFilter("dm_tpl_path", dmTplPath)
	pongo2.RegisterFilter("dm_format_time", dmFormatTime)
	pongo2.RegisterFilter("dm_format_number", dmFormatNumber)
	pongo2.RegisterFilter("dm_wash_number", dmWashNumber)
	pongo2.RegisterFilter("dm_config", dmConfig)
	pongo2.RegisterFilter("dm_json", dmJson)
	pongo2.RegisterFilter("dm_str", dmStr)
	pongo2.RegisterFilter("dm_value", dmValue)
	pongo2.RegisterFilter("dm_now", dmNow)
	pongo2.RegisterFilter("dm_split", dmSplit)
}
