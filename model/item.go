package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct{
	Id      primitive.ObjectID				`bson:"_id,omitempty" json:"item_id"`
	Title       string							`bson:"title" json:"title"`
	Description string							`bson:"description" json:"description"`
	UserId 		primitive.ObjectID				`bson:"userId" json:"userId"`
}



