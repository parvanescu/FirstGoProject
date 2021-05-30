package delivery

import (
	"ExGabi/mutations"
	"ExGabi/payload"
	"ExGabi/response"
	"ExGabi/types"
	"ExGabi/utils/gql"
	"ExGabi/utils/token"
	"fmt"
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
		"updateUserPerformedByLeader": h.updateUserPerformedByLeader(),
		"register": h.register(),
		"addInactiveUser":h.addInactiveUser(),
		"login": h.login(),
		"searchItem":h.searchItem(),
		"checkToken":h.checkToken(),
		"setPassword":h.setPassword(),
		"addPosition":h.addPosition(),
		"exchangePositions":h.exchangePosition(),
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

func (h Handler) updateUserPerformedByLeader() *graphql.Field {
	return &graphql.Field{
		Type: types.UserType,
		Description: "Update a user\n Returns old values of User",
		Args: graphql.FieldConfigArgument{
			"last_name":&graphql.ArgumentConfig{Type: graphql.String},
			"first_name":&graphql.ArgumentConfig{Type: graphql.String},
			"old_email":&graphql.ArgumentConfig{Type: graphql.String},
			"new_email":&graphql.ArgumentConfig{Type: graphql.String},
			"new_password":&graphql.ArgumentConfig{Type: graphql.String},
			"token":&graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams)(interface{},error){
			fmt.Println("aaa")
			responseUsr,err:=h.uC.UpdateUserPerformedByLeader(
				&payload.User{
					LastName:  params.Args["last_name"].(string),
					FirstName: params.Args["first_name"].(string),
					Email:     params.Args["new_email"].(string),
					Password:  params.Args["new_password"].(string),
					Token: params.Args["token"].(string),
				},params.Args["old_email"].(string))
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
			"organisation_name":&graphql.ArgumentConfig{Type: graphql.String},
			"CUI":&graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams)(interface{},error){
			user,organisation,err:=h.uC.Register(
				&payload.User{
					LastName: params.Args["last_name"].(string),
					FirstName: params.Args["first_name"].(string),
					Email:    params.Args["email"].(string),
					Password: "",
				},
				&payload.Organisation{
					Name: params.Args["organisation_name"].(string),
					CUI: params.Args["CUI"].(string),
				})
			user.OrganisationId = organisation.Id
			if err!=nil{
				return nil, err
			}
			return user, nil
		},
	}
}

func (h Handler) addInactiveUser() *graphql.Field {
	return &graphql.Field{
		Type: types.UserType,
		Description: "Add inactive user",
		Args: graphql.FieldConfigArgument{
			"last_name":&graphql.ArgumentConfig{Type: graphql.String},
			"first_name":&graphql.ArgumentConfig{Type: graphql.String},
			"email":&graphql.ArgumentConfig{Type: graphql.String},
			"password":&graphql.ArgumentConfig{Type: graphql.String},
			"token":&graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(params graphql.ResolveParams)(interface{},error){
			user,err:=h.uC.AddInactiveUser(
				&payload.User{
					LastName: params.Args["last_name"].(string),
					FirstName: params.Args["first_name"].(string),
					Email:    params.Args["email"].(string),
					Password: params.Args["password"].(string),
					Token: params.Args["token"].(string),
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
		Type: types.UserType,
		Description: "Login",
		Args: graphql.FieldConfigArgument{
			"email": &graphql.ArgumentConfig{Type: graphql.String},
			"password": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			email := p.Args["email"].(string)
			password := p.Args["password"].(string)
			user,err:=h.uC.Login(&payload.User{Email: email,Password: password})
			return user,err
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

func (h *Handler) setPassword() *graphql.Field {
	return &graphql.Field{
		Type: types.UserType,
		Description: "Sets user password and updates it's status to active",
		Args: graphql.FieldConfigArgument{
			"password": &graphql.ArgumentConfig{Type: graphql.String},
			"user_id": &graphql.ArgumentConfig{Type: gql.ObjectId},
			"organisation_id": &graphql.ArgumentConfig{Type: gql.ObjectId},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			password := p.Args["password"].(string)
			userId := p.Args["user_id"].(primitive.ObjectID)
			organisationId := p.Args["organisation_id"].(primitive.ObjectID)
			err := h.uC.SetUserPassword(&payload.User{Id:userId,Password: password,OrganisationId: organisationId})
			return &response.User{},err
		},
	}
}

func (h Handler) addPosition() *graphql.Field {
	return &graphql.Field{
		Description:       "Add position to a specific organisation",
		Type:              types.PositionType,
		Args:              graphql.FieldConfigArgument{
			"access_level": &graphql.ArgumentConfig{Type: graphql.Int},
			"name": &graphql.ArgumentConfig{Type: graphql.String},
			"token": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			accessLevel := p.Args["access_level"].(int)
			name := p.Args["name"].(string)
			jwtToken := p.Args["token"].(string)
			position,err := h.uC.AddPosition(&payload.Position{
				Name:           name,
				AccessLevel:    accessLevel,
				Token:          jwtToken,
			})
			if err != nil{
				return nil, err
			}
			return position, nil

		},

	}
}

func (h Handler) exchangePosition() *graphql.Field {
	return &graphql.Field{
		Description:       "Exchange positions access levels",
		Type:              graphql.NewList(types.PositionType),
		Args:              graphql.FieldConfigArgument{
			"bottom_position_id": &graphql.ArgumentConfig{Type: gql.ObjectId},
			"bottom_position_level": &graphql.ArgumentConfig{Type: graphql.Int},
			"top_position_id": &graphql.ArgumentConfig{Type: gql.ObjectId},
			"top_position_level": &graphql.ArgumentConfig{Type: graphql.Int},
			"token": &graphql.ArgumentConfig{Type: graphql.String},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			bottomPositionId := p.Args["bottom_position_id"].(primitive.ObjectID)
			bottomPositionLevel := p.Args["bottom_position_level"].(int)
			topPositionId := p.Args["top_position_id"].(primitive.ObjectID)
			topPositionLevel := p.Args["top_position_level"].(int)
			jwtToken := p.Args["token"].(string)
			newBottomPosition,newTopPosition,err := h.uC.ExchangePositionsAccessLevel(
				&payload.Position{Id:bottomPositionId,AccessLevel: bottomPositionLevel,Token: jwtToken},
				&payload.Position{Id:topPositionId,AccessLevel: topPositionLevel})
			if err != nil{
				return nil, err
			}
			positionsList := new([]response.Position)
			positions := append(*positionsList,*newBottomPosition,*newTopPosition)
			return positions, nil

		},

	}
}




