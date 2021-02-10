package types

import (
	"ExGabi/utils/gql"
	"github.com/graphql-go/graphql"
)

var UserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{Type: gql.ObjectId},
			"first_name": &graphql.Field{Type: graphql.String},
			"last_name": &graphql.Field{Type: graphql.String},
			"email": &graphql.Field{Type: graphql.String},
			"password": &graphql.Field{Type: graphql.String},
			"status": &graphql.Field{Type: graphql.Boolean},
			"items": &graphql.Field{Type: graphql.NewList(ItemType)},
			"token": &graphql.Field{Type: graphql.String},
		},
	})
