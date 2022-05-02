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

func main() {
	c := generated.Config{Resolvers: &graph.Resolver{}}
	c.Directives.IsAuth = graph.IsAuth

	router := chi.NewRouter()
	router.Use(auth.Middleware)
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", handler.NewDefaultServer(generated.NewExecutableSchema(c)))

	log.Printf("API is ready")
	panic(http.ListenAndServe(":"+os.Getenv("API_PORT"), router))
}
