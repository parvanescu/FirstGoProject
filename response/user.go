package response

import (
	"ExGabi/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//type User struct{
//	UserId primitive.ObjectID 				`bson:"_id,omitempty" json:"user_id"`
//	UserName string							`bson:"name,omitempty" json:"user_name"`
//}
type User struct{
	model.User
	ItemList []primitive.ObjectID
}

func NewUser(user model.User,itemList []primitive.ObjectID) User{
	return User{
		user,
		itemList,
	}
}

