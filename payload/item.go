package payload

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type Item struct{
	Id      primitive.ObjectID					`bson:"id" json:"id"`
	Title       string							`bson:"title" json:"title"`
	Description string							`bson:"description" json:"description"`
	Token 		string                          `json:"token"`
}

