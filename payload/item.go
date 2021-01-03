package payload

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type Item struct{
	Id      primitive.ObjectID				`bson:"_id" json:"item_id"`
	Title       string							`bson:"title" json:"title"`
	Description string							`bson:"description" json:"description"`
	UserId 		primitive.ObjectID				`bson:"userId" json:"userId"`
	Token 		string                      `bson:"token" json:"token"`
}

