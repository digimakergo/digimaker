//Author xc, Created on 2019-04-28 18:11
//{COPYRIGHTS}

package handler

import (
	"dm/contenttype"
	"dm/fieldtype"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchByContent(t *testing.T) {
	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()

	querier := Querier()
	content, err := querier.FetchByContentID("article", 2)
	assert.Equal(t, nil, err)
	fmt.Println(content.ToMap()["content_id"])
	// assert.Equal(t, 2, content.ToMap()["content_id"].(int))

	content, err = querier.FetchByID(1)
	assert.Equal(t, 1, content.Value("id"))
}
