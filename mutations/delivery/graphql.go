package delivery

import (
	"ExGabi/mutations"
	"ExGabi/payload"
	"ExGabi/response"
	"ExGabi/types"
	"ExGabi/utils/gql"
	"ExGabi/utils/token"
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	uC mutations.IUseCase
}


func New(uc mutations.IUseCase) graphql.Fields {
	h := &Handler{uC: uc}
	return graphql.Fields{
		"addItem": h.addItem(),
		"deleteItem": h.deleteItem(),
		"updateItem": h.updateItem(),
		"deleteUser": h.deleteUser(),
		"updateUser": h.updateUser(),
		"register": h.register(),
		"login": h.login(),
		"searchItem":h.searchItem(),
		"checkToken":h.checkToken(),
	}
}

func (h Handler) addItem() *graphql.Field {
	return &graphql.Field{
		Type: types.ItemType,
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

			responseItem,err:=h.uC.AddItem(&itemToAdd)
			if err!=nil{
				return nil,err
			}
			return responseItem,nil
		},
	}
}

func (h Handler) deleteItem() *graphql.Field {
	return &graphql.Field{
		Type: graphql.String,
		Description: "Delete a note \n Return token",
		Args: graphql.FieldConfigArgument{
			"id":&graphql.ArgumentConfig{
				Type: gql.ObjectId},
			"token":&graphql.ArgumentConfig{
				Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams)(interface{},error){
			id := params.Args["id"].(primitive.ObjectID)
			token:= params.Args["token"].(string)
			itemToDelete:= payload.Item{Id: id,Token: token}
			newTkn,err := h.uC.DeleteItem(&itemToDelete)
			if err!=nil{
				return nil,err
			}
			return newTkn,nil
		},
	}
}

func (h Handler) updateItem() *graphql.Field {
	return &graphql.Field{
		Type: types.ItemType,
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
			responseItem,err:=h.uC.UpdateItem(&itemToUpdate)
			if err!=nil{
				return nil,err
			}
			return responseItem,nil
		},
	}
}

func (h Handler) deleteUser() *graphql.Field {
	return &graphql.Field{
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
			newTkn,err:= h.uC.DeleteUser(&payload.User{Id:id,Token: token})
			if err!=nil{
				return nil,err
			}
			return newTkn,nil
		},
	}
}

func (h Handler) updateUser() *graphql.Field {
	return &graphql.Field{
		Type: types.UserType,
		Description: "Update a user\n Returns old values of User",
		Args: graphql.FieldConfigArgument{
			"last_name":&graphql.ArgumentConfig{Type: graphql.String},
			"first_name":&graphql.ArgumentConfig{Type: graphql.String},
			"email":&graphql.ArgumentConfig{Type: graphql.String},
			"password":&graphql.ArgumentConfig{Type: graphql.String},
			"token":&graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams)(interface{},error){
			responseUsr,err:=h.uC.UpdateUser(
				&payload.User{
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
	}
}

func (h Handler) register() *graphql.Field {
	return &graphql.Field{
		Type: types.UserType,
		Description: "Register user",
		Args: graphql.FieldConfigArgument{
			"last_name":&graphql.ArgumentConfig{Type: graphql.String},
			"first_name":&graphql.ArgumentConfig{Type: graphql.String},
			"email":&graphql.ArgumentConfig{Type: graphql.String},
			"password":&graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams)(interface{},error){
			user,err:=h.uC.Register(
				&payload.User{
					LastName: params.Args["last_name"].(string),
					FirstName: params.Args["first_name"].(string),
					Email:    params.Args["email"].(string),
					Password:  params.Args["password"].(string),
				})
			if err!=nil{
				return nil, err
			}
			return user, nil

		},
	}
}

func (h *Handler) login() *graphql.Field {
	return &graphql.Field{
		Type: graphql.String,
		Description: "Login",
		Args: graphql.FieldConfigArgument{
			"email": &graphql.ArgumentConfig{Type: graphql.String},
			"password": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			email := p.Args["email"].(string)
			password := p.Args["password"].(string)
			tkn,err:=h.uC.Login(&payload.User{Email: email,Password: password})
			return tkn,err
		},
	}
}

func (h *Handler) searchItem() *graphql.Field{
	return &graphql.Field{
		Type: graphql.NewList(types.ItemType),
		Description: "Get items with title or description matching the given key",
		Args: graphql.FieldConfigArgument{
			"criteria": &graphql.ArgumentConfig{Type: graphql.String},
			"token": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			criteria := p.Args["criteria"].(string)
			token := p.Args["token"].(string)
			items,newToken,err:=h.uC.GetMatchingSearch(&payload.Item{Title: criteria,Token: token})


			if len(*items)!=0{
				(*items)[0].Token = newToken
				return items,err
			}else{
				list :=[1]response.Item{}
				list[0]=response.Item{Token: newToken}
				return list,err
			}
		},
	}
}

func (h *Handler)checkToken() *graphql.Field{
	return &graphql.Field{
		Type: graphql.String,
		Description: "Check users token",
		Args: graphql.FieldConfigArgument{
			"token": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			tokenToCheck := p.Args["token"].(string)
			tokenClaims,err := token.CheckToken(tokenToCheck)
			if err!=nil{
				return nil,err
			}
			newToken,err := token.CreateToken(tokenClaims)
			if err!=nil{
				return nil,err
			}
			return newToken,nil
		},
	}
}
