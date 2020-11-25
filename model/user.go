package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct{
	UserId primitive.ObjectID 				`bson:"_id,omitempty" json:"user_id"`
	UserName string							`bson:"name,omitempty" json:"user_name"`
	Password string							`bson:"password,omitempty" json:"password"`
	Status bool                             `bson:"status" json:"status"`
}

func NewUser(name string,password string,status bool) User{
	return User{
		UserName: name,
		Password: password,
		Status: status,
	}
}
