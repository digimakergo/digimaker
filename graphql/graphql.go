package graphql

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/rest"
	"github.com/graphql-go/graphql"
	"github.com/spf13/viper"
)

type postData struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
	Variables map[string]interface{} `json:"variables"`
}

type ArrModel struct {
	ArrInt    []int
	ArrString []string
}

var schema graphql.Schema

func initSchema() {
	schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType(),
			Mutation: mutationType(),
		})
}

func Load() {
	initSchema()

	// try to diff method
	rest.RegisterRoute("/graphql", handleGraphql, "POST")
}

// handleGraphql handles graphQL requests including queries and mutations.
func handleGraphql(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authErr := AuthAPIKey(r)
	if authErr != nil {
		rest.HandleError(authErr, w)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			rest.HandleError(err.(error), w)
			return
		}
	}()

	queryParams := ""
	if http.MethodGet == r.Method {
		queryParams = r.URL.Query().Get("query")
		if queryParams == "" {
			rest.HandleError(errors.New("query is nil"), w)
			return
		}
	} else {
		postParams := postData{}
		if err := json.NewDecoder(r.Body).Decode(&postParams); err != nil {
			rest.HandleError(err, w)
			return
		}
		// set params
		queryParams = postParams.Query
	}
	result := graphql.Do(graphql.Params{
		Context:       ctx,
		Schema:        schema,
		RequestString: queryParams,
	})
	if err := json.NewEncoder(w).Encode(result); err != nil {
		rest.HandleError(err, w)
		return
	}
}

// set graphql/api_key in dm.yaml
func AuthAPIKey(r *http.Request) error {
	apiKey := viper.GetString("graphql.api_key")
	if apiKey == "" {
		log.Error("Not api key set up", "")
		return errors.New("Set up issue")
	}
	rApiKey := r.Header.Get("apiKey")
	if rApiKey == "" {
		return errors.New("Need authorization")
	}
	if rApiKey == apiKey {
		return nil
	} else {
		return errors.New("Wrong api key")
	}
}
