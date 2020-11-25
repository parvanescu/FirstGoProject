package response

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{
	UserId primitive.ObjectID 				`bson:"_id,omitempty" json:"user_id"`
	UserName string							`bson:"name,omitempty" json:"user_name"`
	Status bool                             `bson:"status" json:"status"`
	Items []Item							//`bson:"-" json:"-"`
}

func NewUser(userId primitive.ObjectID,userName string,status bool,items []Item) User{
	return User{
		userId,
		userName,
		status,
		items,
	}
}

