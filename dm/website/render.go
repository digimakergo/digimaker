//Author xc, Created on 2019-06-18 11:20
//{COPYRIGHTS}

//Package website provides website buiding toolkits, template rendering&override. Note: routing is in another package.
package website

import (
	"dm/dm/handler"
	"fmt"
	"net/http"
)

//Given an id RenderContent will output content
func RenderContent(w http.ResponseWriter, r *http.Request, id int) {
	querier := handler.Querier()
	content, err := querier.FetchByID(id)
	fmt.Println(content, err)
}
