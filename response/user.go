package response

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{
	Id primitive.ObjectID 				`bson:"_id,omitempty" json:"user_id"`
	UserName string							`bson:"name,omitempty" json:"user_name"`
	Status bool                             `bson:"status" json:"status"`
	Items []Item							//`bson:"-" json:"-"`
}


