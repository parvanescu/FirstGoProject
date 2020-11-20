package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct{
	UserId primitive.ObjectID 				`bson:"_id,omitempty" json:"user_id"`
	UserName string							`bson:"name,omitempty" json:"user_name"`
}

func NewUser(name string) User{
	return User{
		UserName: name,
	}
}
