package delivery

import (
	"ExGabi/payload"
	"ExGabi/queries"
	"ExGabi/response"
	"ExGabi/types"
	"ExGabi/utils/gql"
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	uC queries.IUseCase
}

func New(uc queries.IUseCase) graphql.Fields {
	h := &Handler{uC: uc}
	return graphql.Fields{
		"getItemById":h.getItemById(),
		"getAllItems":h.getAllItems(),
		"getAllUsersItems":h.getAllUsersItems(),
		"getUserById":h.getUserById(),
		"getUsersProfile":h.getUsersProfile(),
		"getAllUsers":h.getAllUsers(),
		"getAllUsersByOrganisation":h.getAllUsersByOrganisation(),

	}
}

func (h *Handler) getItemById() *graphql.Field {
	return &graphql.Field{
		Type: types.ItemType,
		Description: "Get item by id",
		Args: graphql.FieldConfigArgument{
			"id":&graphql.ArgumentConfig{Type: gql.ObjectId},
			"token":&graphql.ArgumentConfig{
				Type: graphql.String},
		},
		Resolve: func(p graphql.ResolveParams)(interface{},error){
			item,err := h.uC.GetItemById(
				&payload.Item{
					Id:p.Args["id"].(primitive.ObjectID),
					Token: p.Args["token"].(string),
				})
			if err!=nil{
				return nil,err
			}
			return item,err
		},
	}
}

func (h *Handler) getAllItems() *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(types.ItemType),
		Description: "Get all items",
		Args: graphql.FieldConfigArgument{
			"token":&graphql.ArgumentConfig{
				Type: graphql.String},
		},
		Resolve: func(p graphql.ResolveParams)(interface{},error){
			token := p.Args["token"].(string)
			itemList,newToken,err :=h.uC.GetAllItems(token)
			if len(*itemList)!=0{
				(*itemList)[0].Token = newToken
				return itemList,err
			}else{
				list :=[1]response.Item{}
				list[0]=response.Item{Token: newToken}
				return list,err
			}
		},
	}
}

func (h *Handler) getAllUsersItems() * graphql.Field{
	return &graphql.Field{
		Type: graphql.NewList(types.ItemType),
		Description: "Get all items",
		Args: graphql.FieldConfigArgument{
			"token":&graphql.ArgumentConfig{
				Type: graphql.String},
		},
		Resolve: func(p graphql.ResolveParams)(interface{},error){
			token := p.Args["token"].(string)
			itemList,newToken,err :=h.uC.GetAllUsersItems(token)
			if len(*itemList)!=0{
				(*itemList)[0].Token = newToken
				return itemList,err
			}else{
				list :=[1]response.Item{}
				list[0]=response.Item{Token: newToken}
				return list,err
			}
		},
	}
}

func (h *Handler) getUserById() *graphql.Field {
	return &graphql.Field{
		Type: types.UserType,
		Description: "Get user by id",
		Args: graphql.FieldConfigArgument{
			"id":&graphql.ArgumentConfig{Type: gql.ObjectId},
			"token":&graphql.ArgumentConfig{
				Type: graphql.String},
		},
		Resolve: func(p graphql.ResolveParams)(interface{},error){
			id:=p.Args["id"].(primitive.ObjectID)
			token:= p.Args["token"].(string)
			user,err := h.uC.GetUserById(&payload.User{Id:id,Token: token})
			if err!=nil{
				return nil,err
			}
			return user,err
		},
	}
}

func (h *Handler) getUsersProfile() *graphql.Field{
	return &graphql.Field{
		Type: types.UserType,
		Description: "Get logged user's profile ",
		Args: graphql.FieldConfigArgument{
			"token":&graphql.ArgumentConfig{
				Type: graphql.String},
		},
		Resolve: func(p graphql.ResolveParams)(interface{},error){
			token:= p.Args["token"].(string)
			user,err := h.uC.GetUserProfile(token)
			if err!=nil{
				return nil,err
			}
			return user,err
		},
	}
}

func (h *Handler) getAllUsers() *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(types.UserType),
		Description: "Get all users",
		Args: graphql.FieldConfigArgument{
			"token":&graphql.ArgumentConfig{
				Type: graphql.String},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			token:= p.Args["token"].(string)
			usersList,newToken,err:= h.uC.GetAllUsers(token)
			if len(*usersList)!=0{
				(*usersList)[0].Token = newToken
				return usersList,err
			}else{
				list :=[1]response.User{}
				list[0]=response.User{Token: newToken}
				return list,err
			}
		},
	}
}

func (h *Handler) getAllUsersByOrganisation() *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(types.UserType),
		Description: "Get all users by organisation",
		Args: graphql.FieldConfigArgument{
			"token":&graphql.ArgumentConfig{
				Type: graphql.String},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			token:= p.Args["token"].(string)
			usersList,newToken,err:= h.uC.GetAllUsersByOrganisation(token)
			if len(*usersList)!=0{
				(*usersList)[0].Token = newToken
				return usersList,err
			}else{
				list :=[1]response.User{}
				list[0]=response.User{Token: newToken}
				return list,err
			}
		},
	}
}




