package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/olafszymanski/user-cms/auth"
	"github.com/olafszymanski/user-cms/graph"
	"github.com/olafszymanski/user-cms/graph/generated"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	c := generated.Config{Resolvers: &graph.Resolver{}}
	c.Directives.IsAuth = graph.IsAuth

	router := chi.NewRouter()
	router.Use(auth.Middleware)
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", handler.NewDefaultServer(generated.NewExecutableSchema(c)))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
