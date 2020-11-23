package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserItem struct{
	UserId primitive.ObjectID				`bson:"userId" json:"userId"`
	ItemId primitive.ObjectID				`bson:"itemId" json:"itemId"`
}

func NewUserItem(newUserId primitive.ObjectID,newItemId primitive.ObjectID)UserItem{
	return UserItem{newUserId,newItemId}
}
