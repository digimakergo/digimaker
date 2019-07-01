//Author xc, Created on 2019-06-18 11:20
//{COPYRIGHTS}

//Package website provides website buiding toolkits, template rendering&override. Note: routing is in another package.
package sitekit

//
// //Given an id RenderContent will output content
// func RenderContent(context context.Context, id int, userID int) {
// 	querier := handler.Querier()
// 	content, err := querier.FetchByID(id)
// 	variables := map[string]interface{}{}
// 	if err == nil {
// 		variables["error"] = "1101"
// 	} else {
// 		location := content.GetLocation()
// 		canRead, err := handler.CanRead(userID, content, context)
// 		if err != nil {
// 			variables["error"] = "1100"
// 		}
//
// 		if !canRead {
// 			variables["error"] = "1103"
// 		}
//
// 		if location.IsInvisible {
// 			variables["error"] = "1104" //todo: more on hidden content based on user.
// 		}
// 	}
//
// 	variables["content"] = content
// 	return variables
// }
