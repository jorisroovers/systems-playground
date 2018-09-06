package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
)

type GraphQLHttpQuery struct {
	Query string `json:"query"`
}

func main() {
	// define simple graphQL Schema: { hello }
	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world (custom handler)", nil
			},
		},
	}
	// Setup schema for use
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Standard go HTTP handler
	// Takes in a json query of format: { "query": "{ hello }" } and the query on to the GraphQL backend
	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		// Parse incoming json from request body
		decoder := json.NewDecoder(req.Body)
		data := GraphQLHttpQuery{}
		err := decoder.Decode(&data)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Incoming graphQL query on /hello: ", data.Query)

		// Execute the graphQL query
		params := graphql.Params{Schema: schema, RequestString: data.Query}
		res := graphql.Do(params)
		if len(res.Errors) > 0 {
			log.Fatalf("failed to execute graphql operation, errors: %+v", res.Errors)
		}

		// Marshall response to json and send it back
		rJSON, _ := json.Marshal(res)
		fmt.Fprintf(w, string(rJSON))
		log.Printf("Response: %s \n", rJSON) // {“data”:{“hello”:”world”}}

	})
	log.Println("Starting server...")
	http.ListenAndServe(":5678", nil)
}
