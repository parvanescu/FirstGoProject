package response

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{
	Id primitive.ObjectID 				`bson:"_id,omitempty" json:"user_id"`
	LastName string							`bson:"lastName,omitempty" json:"last_name"`
	FirstName string							`bson:"firstName,omitempty" json:"first_name"`
	Email string 							`bson:"email" json:"email"`
	Status bool                             `bson:"status" json:"status"`
	Items []Item							`bson:"items" json:"items"`
	Token string							`json:"token"`
	jwt.StandardClaims
}


