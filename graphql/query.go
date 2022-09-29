package graphql

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/graphql-go/graphql"

	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/query"
	"github.com/digimakergo/digimaker/rest"
	"github.com/gorilla/mux"
)

func GraphqlList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cType := "article"
	queryParams := r.URL.Query().Get("query")
	if queryParams == "" {
		rest.HandleError(errors.New("query is nil"), w)
		return
	}

	fmt.Println("queryParams:", queryParams)

	list, _, err := query.List(ctx, cType, db.EmptyCond().Limit(0, 10))
	if err != nil {
		rest.HandleError(err, w)
		return
	}
	var listMap []map[string]interface{}
	for _, item := range list {
		m, _ := contenttype.ContentToMap(item)
		listMap = append(listMap, m)
	}

	// generate graphql type
	typeModel := contenttype.NewInstance(cType)
	retMaps, _ := contenttype.ContentToMap(typeModel)

	fields := graphql.Fields{}
	for key, value := range retMaps {
		switch value.(type) {
		case string:
			fields[key] = &graphql.Field{Type: graphql.String}
		case int:
			fields[key] = &graphql.Field{Type: graphql.Int}
		case interface{}:
			fields[key] = &graphql.Field{Type: &graphql.Interface{}}
		}
	}

	var modelType = graphql.NewObject(
		graphql.ObjectConfig{
			Name:   cType,
			Fields: fields,
		})

	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"list": &graphql.Field{
					Type:        graphql.NewList(modelType),
					Description: "get list by type",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						fmt.Println(p.Info.Schema)
						return listMap, nil
					},
				},
			},
		})

	var schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryType,
		})

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: queryParams,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors : %v", result.Errors)
	}

	rest.WriteResponse(result, w)
}

func init() {
	rest.RegisterRoute("/graphql", GraphqlList, "POST")
}
