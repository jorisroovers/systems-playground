# Graphql-CouchDB
Simple examples of using [GraphQL](https://graphql.org/) in Go.

As often, there's a bunch of competing GraphQL libraries for Go out there.
- [graphql-go/graphql](https://github.com/graphql-go/graphql)
- [graph-gophers/graphql-go](https://github.com/graph-gophers/graphql-go)
- [99designs/gqlgen](https://github.com/99designs/gqlgen)

For this repo, we're sticking with **graphql-go/graphql**. It seems like the others have benefits over graphql-go/graphql (in particular around ease of creating schema definitions), but graphql-go is the most popular one on github. The goal of this playground is not finding the best library, yet just learning about GraphQL (and Go) :-)

On the front-end, apollo seems like a good library: https://www.apollographql.com/client/

## GraphQL Basics

### GraphQL queries

You can play around on [fakerql.com](https://fakerql.com/) with some graphQL queries.
The rest of the notes in this section can just be copy-pasted in there:

```graphql
# GraphQL allows comments with #. Whitespace is just for formatting, has no effect
# Get all users, and get id and lastname for each user
{
    allUsers {
        id
        lastName
    }
}

# All requests send to graphQL are so-called 'operations'. The default is 'query', another common one is 'mutation'
# Since 'query' is the default, you don't have the specify it, but you *can*, like so:
query GetAllUsers{
    allUsers {
        id
        lastName
    }
}

# You can even specificy your own name for the query
query MyCustomName{
    allUsers {
        id
        lastName
    }
}

# Note that all of this can also be done from the command-line:
# curl -s -H "content-type: application/graphql" -d "{allUsers {id}}" https://fakerql.com/graphql | jq .
# Same query but using json body
# curl -s -H "content-type: application/json" -d '{"query": " { allUsers {id}}"}' https://fakerql.com/graphql | jq .

# Get Todos:
query AllTodos {
  allTodos {
    id
    title
  }
}

# Get a specific TODO
{
  Todo(id:"cjlqqeamr152y3g10l8mh656a"){
    id
    title
    completed
  }
}

# You can use mutations to make modifications to datasets:
# Note: This mutation operation is accepted, but doesn't acctually mutate the dataset on fakerql. I believe this is just
# the backend not persisting changes, not an issue with the mutation itself.
mutation AddTodo {
  createTodo(title:"foobar" , completed: false) {
    id
    title
    completed
  }
}

# You can go meta and query the schema itself, like so:
# Describe available fields for type Post
{ __type(name:"Post") {
    fields {
      name
      description
      }
    }
}
```

GraphQL is a lot more powerful than these simple examples, for more detail, refer to the [GraphQL website](http://graphql.github.io/learn/).

### GraphQL over JSON

Typically, GraphQL is transported over JSON using the following format.
Request use the "query" key:
```json
{ "query": "{ hello }" }
```
Responses have a "data" key:
```json
{
    "data": {
        "hello" : "world" 
    }
}
```

## Programs

There's multiple programs in this directory, listed here in order of increasing complexity.

```bash
go get
```

### graphql-simple.go

Uses [graphql-go/graphql](http://github.com/graphql-go/graphql) to create a graphQL schema and serves it over HTTP using
[graphql-go/handler](http://github.com/graphql-go/handler).

The goal of this script is to showcase the basics of GraphQL.
```bash
# Run server
go run graphql-simple.go

# Client request (in different terminal session)
# The server can take the client request in different formats based on the header.
# content-type: application/json
curl -H 'content-type: application/json' -d '{"query":"{ hello }"}'  "localhost:5678/hello"
# content-type: application/graphql
curl -H 'content-type: application/graphql' -d 'query { hello }'  "localhost:5678/hello"
# No header (=use query string, note escaping curly brackets for use in shell)
curl -v "localhost:5678/hello?query=\{hello\}"
```

#### GraphiQL
The handler also has [GraphiQL](https://github.com/graphql/graphiql) enabled (a inbrowser GraphQL IDE), so if you browse
to ```localhost:5678/hello```, you'll see that pop up.
You can also query the schema, e.g.:
```graphql
{
  __schema {
    queryType {
      name
      description
      fields{
        name
      }
    }
  }
}
```

### graphql-custom-handler.go

Uses [graphql-go/graphql](http://github.com/graphql-go/graphql) to create a graphQL schema and then manually wires this
to a HTTP server using Go's built-in ```net/http```.

 This program showcases that GraphQL doesn't necessarily need to run over HTTP or use JSON, that part is independent from the core GraphQL engine and implemented manually in this example.

```bash
# Run server
go run graphql-custom-handler.go

# Client request (in different terminal session)
curl -d '{"query":"{ hello }"}'  "localhost:5678/hello"
```

### graphql-complex-data.go


```bash
# Run server
go run graphql-complex-data.go
```
Then browse to ```localhost:5678/hello``` and use the GraphiQL IDE to enter the following queries

```graphql
# specific address
{
  address(street: "Wall Street", number: 2001) {
    street
    number
    country
  }
}

# Multiple addresses, more fields
{
  addresses {
    street
    number
    city
    fullAddress
  }
}


# List all people
# Note: It's not possible to filter on people -> not implemented in this example for type Person
{
  people {
    name
    address {
      country
    }
  }
}

# Filter a specific person
{
  person(name: "John") {
    name
    address {
      street
      country
      number
    }
  }
}
```