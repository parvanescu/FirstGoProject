package main

import (
	"ExGabi/payload"
	"ExGabi/repository"
	"ExGabi/response"
	"ExGabi/useCase"
	"ExGabi/utils/gql"
	"context"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	var repo repository.IRepository = repository.New(client)
	var uC useCase.IUseCase = useCase.New(repo)


	var itemType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Item",
			Fields: graphql.Fields{
				"id": &graphql.Field{Type: gql.ObjectId},
				"title": &graphql.Field{Type: graphql.String},
				"description": &graphql.Field{Type: graphql.String},
				"userId":&graphql.Field{Type: gql.ObjectId},
				"token": &graphql.Field{Type: graphql.String}},
		},
		)

	var userType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "User",
			Fields: graphql.Fields{
				"id": &graphql.Field{Type: gql.ObjectId},
				"first_name": &graphql.Field{Type: graphql.String},
				"last_name": &graphql.Field{Type: graphql.String},
				"email": &graphql.Field{Type: graphql.String},
				"password": &graphql.Field{Type: graphql.String},
				"status": &graphql.Field{Type: graphql.Boolean},
				"items": &graphql.Field{Type: graphql.NewList(itemType)},
				"token": &graphql.Field{Type: graphql.String},
			},
		})



	var mutationType = graphql.NewObject(graphql.ObjectConfig{
		Name:"Mutation",
		Fields: graphql.Fields{
			"addItem": &graphql.Field{
				Type: itemType,
				Description: "Add new note.\nReturns an item with id and token field set.\nThe token is refreshed.",
				Args: graphql.FieldConfigArgument{
					"title":&graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String) },
					"description":&graphql.ArgumentConfig{
							Type: graphql.String},
					"token":&graphql.ArgumentConfig{
							Type: graphql.String},
				},
				Resolve: func(params graphql.ResolveParams)(interface{},error){
					token:= params.Args["token"].(string)
					itemToAdd:= payload.Item{
						Title: params.Args["title"].(string),
						Description: params.Args["description"].(string),
						Token: token}

					responseItem,err:=uC.AddItem(&itemToAdd)
					if err!=nil{
						return nil,err
					}
					return responseItem,nil
				},
			},
			"deleteItem":&graphql.Field{
				Type: graphql.String,
				Description: "Delete a note \n Return token",
				Args: graphql.FieldConfigArgument{
					"id":&graphql.ArgumentConfig{
						 Type: graphql.NewNonNull(gql.ObjectId)},
					"token":&graphql.ArgumentConfig{
						 	Type: graphql.NewNonNull(graphql.String)},
						 	},
				Resolve: func(params graphql.ResolveParams)(interface{},error){
					id := params.Args["id"].(primitive.ObjectID)
					token:= params.Args["token"].(string)
					itemToDelete:= payload.Item{Id: id,Token: token}
					newTkn,err := uC.DeleteItem(&itemToDelete)
					if err!=nil{
						return nil,err
					}
					return newTkn,nil
				},
			},
			"updateItem":&graphql.Field{
				Type: itemType,
				Description: "Update a note\n Returns the old item and the refreshed token in the same object\n",
				Args: graphql.FieldConfigArgument{
					"id":&graphql.ArgumentConfig{
						Type: graphql.NewNonNull(gql.ObjectId)},
					"title":&graphql.ArgumentConfig{
						Type: graphql.String },
					"description":&graphql.ArgumentConfig{
						Type: graphql.String},
					"token":&graphql.ArgumentConfig{
						Type: graphql.String},
					},
				Resolve: func(params graphql.ResolveParams)(interface{},error){
					itemToUpdate:= payload.Item{
						Id: params.Args["id"].(primitive.ObjectID),
						Title: params.Args["title"].(string),
						Description:  params.Args["description"].(string),
						Token: params.Args["token"].(string),
					}
					responseItem,err:=uC.UpdateItem(&itemToUpdate)
					if err!=nil{
						return nil,err
					}
					return responseItem,nil
				},
			},
			"deleteUser":&graphql.Field{
				Type: graphql.String,
				Description: "Delete a user \n Return token",
				Args: graphql.FieldConfigArgument{
					"id":&graphql.ArgumentConfig{Type: gql.ObjectId},
					"email":&graphql.ArgumentConfig{Type: graphql.String},
					"token":&graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(params graphql.ResolveParams)(interface{},error){
					id:=params.Args["id"].(primitive.ObjectID)
					token:=params.Args["token"].(string)
					newTkn,err:= uC.DeleteUser(&payload.User{Id:id,Token: token})
					if err!=nil{
						return nil,err
					}
					return newTkn,nil
				},
			},
			"updateUser":&graphql.Field{
				Type: userType,
				Description: "Update a user\n Returns old values of User",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: gql.ObjectId},
					"last_name":&graphql.ArgumentConfig{Type: graphql.String},
					"first_name":&graphql.ArgumentConfig{Type: graphql.String},
					"email":&graphql.ArgumentConfig{Type: graphql.String},
					"password":&graphql.ArgumentConfig{Type: graphql.String},
					"token":&graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(params graphql.ResolveParams)(interface{},error){
					responseUsr,err:=uC.UpdateUser(
						&payload.User{
						Id: params.Args["id"].(primitive.ObjectID),
						LastName:  params.Args["last_name"].(string),
						FirstName: params.Args["first_name"].(string),
						Email:     params.Args["email"].(string),
						Password:  params.Args["password"].(string),
						Token: params.Args["token"].(string),
					})
					if err!=nil{
						return nil, err
					}
					return responseUsr, nil

				},
			},
			"register":&graphql.Field{
				Type: graphql.String,
				Description: "Register user",
				Args: graphql.FieldConfigArgument{
					"last_name":&graphql.ArgumentConfig{Type: graphql.String},
					"first_name":&graphql.ArgumentConfig{Type: graphql.String},
					"email":&graphql.ArgumentConfig{Type: graphql.String},
					"password":&graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(params graphql.ResolveParams)(interface{},error){
					tkn,err:=uC.Register(
						&payload.User{
						LastName: params.Args["last_name"].(string),
						FirstName: params.Args["first_name"].(string),
						Email:    params.Args["email"].(string),
						Password:  params.Args["password"].(string),
						})
					if err!=nil{
						return nil, err
					}
					return tkn, nil

				},
			},
		},
	})

	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name:"Query",
			Fields: graphql.Fields{
				"getItemById": &graphql.Field{
					Type: itemType,
					Description: "Get item by id",
					Args: graphql.FieldConfigArgument{
						"id":&graphql.ArgumentConfig{Type: gql.ObjectId},
						"token":&graphql.ArgumentConfig{
							Type: graphql.String},
					},
					Resolve: func(p graphql.ResolveParams)(interface{},error){
						if err!=nil{
							return nil,err
						}
						item,err := uC.GetItemById(
							&payload.Item{
								Id:p.Args["id"].(primitive.ObjectID),
								Token: p.Args["token"].(string),
							})
						return item,err
					},
				},
				"getAllItems": &graphql.Field{
					Type: graphql.NewList(itemType),
					Description: "Get all items",
					Args: graphql.FieldConfigArgument{
						"token":&graphql.ArgumentConfig{
							Type: graphql.String},
					},
					Resolve: func(p graphql.ResolveParams)(interface{},error){
						token := p.Args["token"].(string)
						itemList,newToken,err :=uC.GetAllItems(token)
						if len(*itemList)!=0{
							(*itemList)[0].Token = newToken
							return itemList,err
						}else{
							list :=[1]response.Item{}
							list[0]=response.Item{Token: newToken}
							return list,err
						}
					},
				},
				"getUserById": &graphql.Field{
					Type: userType,
					Description: "Get user by id",
					Args: graphql.FieldConfigArgument{
						"id":&graphql.ArgumentConfig{Type: gql.ObjectId},
						"token":&graphql.ArgumentConfig{
							Type: graphql.String},
					},
					Resolve: func(p graphql.ResolveParams)(interface{},error){
						id:=p.Args["id"].(primitive.ObjectID)
						token:= p.Args["token"].(string)
						if err!=nil{
							return nil,err
						}
						user,err := uC.GetUserById(&payload.User{Id:id,Token: token})
						return user,err
					},
				},
				"getAllUsers": &graphql.Field{
					Type: graphql.NewList(userType),
					Description: "Get all users",
					Args: graphql.FieldConfigArgument{
						"token":&graphql.ArgumentConfig{
							Type: graphql.String},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						token:= p.Args["token"].(string)
						usersList,newToken,err:= uC.GetAllUsers(token)
						if len(*usersList)!=0{
							(*usersList)[0].Token = newToken
							return usersList,err
						}else{
							list :=[1]response.User{}
							list[0]=response.User{Token: newToken}
							return list,err
						}
					},
				},
				"login": &graphql.Field{
					Type: graphql.String,
					Description: "Login",
					Args: graphql.FieldConfigArgument{
						"email": &graphql.ArgumentConfig{Type: graphql.String},
						"password": &graphql.ArgumentConfig{Type: graphql.String},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						email := p.Args["email"].(string)
						password := p.Args["password"].(string)
						tkn,err:=uC.Login(&payload.User{Email: email,Password: password})
						return tkn,err
					},
				},
			},
		})


		var schema, _ = graphql.NewSchema(
			graphql.SchemaConfig{
				Query: queryType,
				Mutation: mutationType,
			})



		hGraphQl := InjectGraphqlHandler(&schema)

		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedHeaders:   []string{"Access-Control-Allow-Headers", "Origin", "Accept", "Content-Type", "Content-Length", "Access-Control-Request-Method", "Access-Control-Request-Headers", "X-Requested-With", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
			AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowCredentials: true,
			Debug:            false,
		})

		mux := http.NewServeMux()
		mux.HandleFunc("/graphql", CorsMiddleware(hGraphQl))

		h := c.Handler(mux)
		log.Printf("Serving on http://localhost:8080")
		log.Panicln("HTTP server error: ", http.ListenAndServe(":8080", h))


}

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
