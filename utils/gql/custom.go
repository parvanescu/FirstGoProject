package gql

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ObjectId = graphql.NewScalar(graphql.ScalarConfig{
	Name: "BSON",
	Description: "A BSON OBJECT",
	Serialize: func(value interface{})interface{}{
		switch value:=value.(type) {
		case primitive.ObjectID:
			return value.Hex()
		case *primitive.ObjectID:
			v:=*value
			return v.Hex()
		default:
			return nil
		}
	},
	ParseValue: func(value interface{})interface{}{
		switch value:=value.(type){
		case string:
			id,_:=primitive.ObjectIDFromHex(value)
			return id
		case *string:
			id,_:=primitive.ObjectIDFromHex(*value)
			return id
		default:
			return nil
		}
	},
	ParseLiteral: func(valueAST ast.Value)interface{}{
		switch valueAST.(type) {
		case *ast.StringValue:
			id,_:=primitive.ObjectIDFromHex(valueAST.GetValue().(string))
			return id
		}
		return nil
	},
})

