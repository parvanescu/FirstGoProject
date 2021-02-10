package main

import (
	"ExGabi/utils/server"
	"context"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)


func main(){
	uri := "mongodb+srv://DbUser123:Password123@cluster0.5tzzs.mongodb.net/<dbname>?retryWrites=true&w=majority"
	ctx,cancel := context.WithTimeout(context.TODO(),10*time.Second)
	defer cancel()
	client,err := mongo.Connect(ctx,options.Client().ApplyURI(uri))
	if err != nil{
		panic(err)
	}


	graphSchema, err := server.InjectSchema(server.MutationInit(client), server.QueryInit(client))
	if err != nil {
		panic(err)
	}
	hGraphQL := server.InjectGraphqlHandler(graphSchema)

		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedHeaders:   []string{"Access-Control-Allow-Headers", "Origin", "Accept", "Content-Type", "Content-Length", "Access-Control-Request-Method", "Access-Control-Request-Headers", "X-Requested-With", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
			AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowCredentials: true,
			Debug:            false,
		})

		mux := http.NewServeMux()
		mux.HandleFunc("/graphql", server.CorsMiddleware(hGraphQL))

		h := c.Handler(mux)
		log.Printf("Serving on http://localhost:8080")
		log.Panicln("HTTP server error: ", http.ListenAndServe(":8080", h))


}
