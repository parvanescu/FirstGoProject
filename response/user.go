package response

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct{
	UserId primitive.ObjectID 				`bson:"_id,omitempty" json:"user_id"`
	UserName string							`bson:"name,omitempty" json:"user_name"`
	PostedItemIds []primitive.ObjectID 		`bson:"itemsIds,omitempty" json:"posted_item_ids"`
}

func NewUser(name string) User{
	return User{
		UserName: name,
		PostedItemIds: []primitive.ObjectID{},
	}
}

