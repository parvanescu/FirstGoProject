package server

import (
	mutationDelivery "ExGabi/mutations/delivery"
	mutationRepository "ExGabi/mutations/repository"
	mutationUseCase "ExGabi/mutations/useCase"
	queriesDelivery "ExGabi/queries/delivery"
	queriesRepository "ExGabi/queries/repository"
	queriesUseCase "ExGabi/queries/useCase"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func InjectGraphqlHandler(schema *graphql.Schema) *handler.Handler {
	h := handler.New(&handler.Config{
		Schema:     schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})
	return h
}

func CorsMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, "+
			"X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func MutationInit(c *mongo.Client) graphql.Fields {
	repo := mutationRepository.New(c)
	useCase := mutationUseCase.New(repo)
	return mutationDelivery.New(useCase)
}

func QueryInit(c *mongo.Client) graphql.Fields {
	repo := queriesRepository.New(c)
	useCase := queriesUseCase.New(repo)
	return queriesDelivery.New(useCase)
}


func InjectSchema(m, q graphql.Fields) (*graphql.Schema, error) {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: q,
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Mutations",
			Fields: m,
		}),
	})
	return &schema, err
}
