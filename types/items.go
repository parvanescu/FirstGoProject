package types

import (
	"ExGabi/utils/gql"
	"github.com/graphql-go/graphql"
)

var ItemType = graphql.NewObject(
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
