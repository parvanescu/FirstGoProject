package payload

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{
	Id primitive.ObjectID 				`bson:"_id,omitempty" json:"user_id"`
	LastName string							`bson:"lastName,omitempty" json:"last_name"`
	FirstName string							`bson:"firstName,omitempty" json:"first_name"`
	Email string 							`bson:"email" json:"email"`
	Password string							`bson:"password,omitempty" json:"password"`
	Status bool                             `bson:"status" json:"status"`
	Token string                    `bson:"token" json:"token"`
}
