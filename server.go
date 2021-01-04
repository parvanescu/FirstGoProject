package main

import (
	"ExGabi/graph"
	"ExGabi/graph/generated"
	"ExGabi/repository"
	"ExGabi/useCase"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	uri := "mongodb+srv://DbUser123:Password123@cluster0.5tzzs.mongodb.net/<dbname>?retryWrites=true&w=majority"
	ctx,cancel := context.WithTimeout(context.TODO(),10*time.Second)
	defer cancel()
	client,err := mongo.Connect(ctx,options.Client().ApplyURI(uri))
	if err != nil{
		panic(err)
	}
	var repo repository.IRepository = repository.New(client)
	var uC useCase.IUseCase = useCase.New(repo)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{uC}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
