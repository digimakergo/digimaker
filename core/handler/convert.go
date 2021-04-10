package handler

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/query"
)

//ConvertToHtml converts html which has content information to real html
func ConvertToHtml(ctx context.Context, input string, removeDataAttribute bool, imagePrefix string) string {

	//convert image to updated image path
	re := regexp.MustCompile(`<img[^>]+data-dm-content="[^"]+"[^>]+>`)

	replaceFunc := func(currentStr string) string {
		re2 := regexp.MustCompile(`([^ =]+)="([0-9a-zA-Z]|;)+"`)
		attributes := re2.FindAllString(currentStr, -1)
		attributeMap := map[string]string{}
		for _, attStr := range attributes {
			arr := strings.Split(attStr, "=")
			name := arr[0]
			value := strings.ReplaceAll(arr[1], `"`, "")
			attributeMap[name] = value
		}
		contentInfo := strings.Split(attributeMap["data-dm-content"], ";")
		if len(contentInfo) <= 1 {
			log.Warning("data-dm-content has wrong format, should be <contenttype>;<cuid>, no replace done. - "+currentStr, "output", ctx)
			return currentStr
		}

		content, _ := query.FetchByCUID(context.Background(), contentInfo[0], contentInfo[1])
		widthStr := ""
		if width, ok := attributeMap["width"]; ok {
			widthStr = `width="` + width + `"`
		}

		heightStr := ""
		if height, ok := attributeMap["height"]; ok {
			heightStr = `height="` + height + `"`
		}

		dataAttribute := ""
		if !removeDataAttribute {
			dataAttribute = `data-dm-content="` + attributeMap["data-dm-content"] + `"`
		}

		if content == nil {
			//to do: check reason(might be missing access) and give log, and output different image
			return fmt.Sprintf(`<img src="not-available.png" %v %v %v />`, widthStr, heightStr, dataAttribute) //todo: make it configurable
		}

		path := imagePrefix + content.Value("image").(string)

		result := fmt.Sprintf(`<img src="%v" %v %v %v />`, path, widthStr, heightStr, dataAttribute)
		return result
	}
	result := re.ReplaceAllStringFunc(input, replaceFunc)
	return result
}
