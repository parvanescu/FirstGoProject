package types

import (
	"ExGabi/utils/gql"
	"github.com/graphql-go/graphql"
)

var PositionType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Position",
		Fields: graphql.Fields{
			"id": &graphql.Field{Type: gql.ObjectId},
			"organisation_id": &graphql.Field{Type: gql.ObjectId},
			"name": &graphql.Field{Type: graphql.String},
			"access_level":&graphql.Field{Type: graphql.Int},
			"token": &graphql.Field{Type: graphql.String}},
	},
)
