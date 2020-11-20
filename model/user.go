package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct{
	UserId primitive.ObjectID `bson:"_id,omitempty"`
	UserName string`bson:"name,omitempty"`
	PostedItemIds []primitive.ObjectID `bson:"itemsIds,omitempty"`
}

func NewUser(name string) User{
	return User{
		UserName: name,
		PostedItemIds: []primitive.ObjectID{},
	}
}
