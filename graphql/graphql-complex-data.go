package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

/*
	# This is the schema we'll be defining:
	type Person {
		name: String
		gender: Gender
		address: Address
	}
	enum Gender {
		MALE
		FEMALE
		OTHER
	}
	type Address {
		street: String! # excl. mark means it's non-nullable: will always return a value
		number: int
		numberExtra: String #e.g. the 'A' in 112A
		city: String
		province: String
		country: String
	}
*/

//// DATA TYPES  ///////////////////////////////////////////////////////////////////////////////////////////////////////
// Data structs used to store our dummy data. Note that this has nothing to do with GraphQL specifically, we just use
// this for some of the wiring to make a more realistic example.

type Address struct {
	Street   string
	Number   int
	City     string
	Province string
	Country  string
}

type Person struct {
	Name    string
	Address Address
	// there's no real enum type in Go, let's ignore that for the purpose of this program :-)
	Gender int
}

func main() {

	//// DATA  /////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Dataset: this is the data that we'll return later
	// Normally, this is stored in e.g. a database, but we're keeping it simple here
	addresses := []Address{
		Address{"Wall Street", 2001, "New York", "New York", "USA"}, Address{"Wall Street", 9123, "New York", "New York", "USA"},
		Address{"Main Street", 5, "San Francisco", "California", "USA"}, Address{"Main Street", 100, "San Francisco", "California", "USA"},
		Address{"Kerk Straat", 13, "Bavel", "Noord-Brabant", "Netherlands"},
	}

	people := []Person{
		Person{"John", addresses[0], 1}, Person{"Jane", addresses[1], 2}, Person{"Maria", addresses[2], 1},
		Person{"Mark", addresses[3], 3}, Person{"Sophia", addresses[4], 2}, Person{"Peter", addresses[0], 3},
	}

	//// GRAPHQL SCHEMA ////////////////////////////////////////////////////////////////////////////////////////////////

	genderEnum := graphql.NewEnum(graphql.EnumConfig{
		Name: "Gender",
		Values: graphql.EnumValueConfigMap{
			"Male": &graphql.EnumValueConfig{
				Value: 1,
			},
			"Female": &graphql.EnumValueConfig{
				Value: 2,
			},
			"Other": &graphql.EnumValueConfig{
				Value: 3,
			},
		},
	})

	addressType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Address",
		Fields: graphql.Fields{
			"street": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String), // recall, we chose street to be mandatory (enforced by NewNonNull(...))
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if address, ok := p.Source.(Address); ok {
						return address.Street, nil
					}
					// recall, we can't return Nil because Street is non-nullable (maybe not the best example...)
					return "NOT FOUND", nil
				},
			},
			"number": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if address, ok := p.Source.(Address); ok {
						return address.Number, nil
					}
					return nil, nil
				},
			},
			"city": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if address, ok := p.Source.(Address); ok {
						return address.City, nil
					}
					return nil, nil
				},
			},
			"province": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if address, ok := p.Source.(Address); ok {
						return address.Province, nil
					}
					return nil, nil
				},
			},
			"country": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if address, ok := p.Source.(Address); ok {
						return address.Country, nil
					}
					return nil, nil
				},
			},
			// This is a custom field that doesn't map directly on a field in the struct, but returns a combination
			// of fields
			"fullAddress": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if address, ok := p.Source.(Address); ok {
						return strconv.Itoa(address.Number) + ", " + address.Street + "\n" +
							address.City + ", " + address.Province + "\n" +
							address.Country, nil
					}
					return nil, nil
				},
			},
		},
	})

	personType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Person",
		Description: "A person", // Description is optional, just using it once here to point that out
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if person, ok := p.Source.(Person); ok {
						return person.Name, nil
					}
					return nil, nil
				},
			},
			"address": &graphql.Field{
				Type: addressType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if person, ok := p.Source.(Person); ok {
						return person.Address, nil
					}
					return nil, nil
				},
			},
			"gender": &graphql.Field{
				Type: genderEnum,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if person, ok := p.Source.(Person); ok {
						return person.Gender, nil
					}
					return nil, nil
				},
			},
		},
	})

	// The root fields here are those that are actually exposed to the end-user as query endpoints
	rootFields := graphql.Fields{
		// people endpoint lists all person(s). Note that you cannot filter on people, see 'addresses' endpoint for that
		// below
		"people": &graphql.Field{
			Type: graphql.NewList(personType), // Note the Type: NewList(...) because we always return a list
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return people, nil
			},
		},
		"person": &graphql.Field{
			Type: personType,
			// Args make it possible to do filtering (or basically whatever you want with it in the Resolve function,
			// but most often used for filtering)
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Description: "Name of the Person",
					Type:        graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Specific street requested, return those street names that match
				for _, person := range people {
					if person.Name == p.Args["name"] {
						return person, nil
					}
				}
				return nil, nil
			},
		},
		// Allow searching addresses on multiple Args
		"address": &graphql.Field{
			Type: addressType,
			Args: graphql.FieldConfigArgument{
				"street": &graphql.ArgumentConfig{
					Description: "Name of the street",
					Type:        graphql.NewNonNull(graphql.String),
				},
				"number": &graphql.ArgumentConfig{
					Description: "Number of the street",
					Type:        graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Specific street requested, return those street names that match
				for _, address := range addresses {
					if address.Street == p.Args["street"] && p.Args["number"] == address.Number {
						return address, nil
					}
				}
				return nil, nil
			},
		},
		// Query endpoint to search multiple streets, type = NewList(...)
		"addresses": &graphql.Field{
			Type: graphql.NewList(addressType),
			Args: graphql.FieldConfigArgument{
				"street": &graphql.ArgumentConfig{
					Description: "Name of the street",
					Type:        graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if p.Args["street"] != nil {
					// Specific street requested, return those street names that match
					matches := []Address{}
					for _, address := range addresses {
						if address.Street == p.Args["street"] {
							matches = append(matches, address)
						}
					}
					return matches, nil
				}
				// else: Return all addresses
				return addresses, nil
			},
		},
	}

	rootQuery := graphql.NewObject(graphql.ObjectConfig{Name: "RootQuery", Fields: rootFields})
	// Setup schema for use, note that how we specify 'Query' here. We could also specify 'Mutation' here if we wanted to
	// support that. Implementing mutations is very similar to Queries.
	schemaConfig := graphql.SchemaConfig{Query: rootQuery}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Use library http handler to expose GraphQL over HTTP
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// Actual HTTP server is Go's built-in
	http.Handle("/hello", h)
	http.ListenAndServe(":5678", nil)

}
