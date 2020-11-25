package payload

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct{
	ItemId      primitive.ObjectID				`bson:"_id" json:"item_id"`
	Title       string							`bson:"title" json:"title"`
	Description string							`bson:"description" json:"description"`
	UserId 		primitive.ObjectID				`bson:"userId" json:"userId"`
}

func NewItem(title string,description string,userId primitive.ObjectID) Item {
	return Item{
		Title:       title,
		Description: description,
		UserId: userId,
	}
}
