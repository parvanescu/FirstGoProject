package response

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct{
	ItemId      primitive.ObjectID				`bson:"_id,omitempty" json:"_id"`
	Title       string							`bson:"title,omitempty" json:"title"`
	Description string							`bson:"description,omitempty" json:"description"`
	UserId 		primitive.ObjectID				`bson:"userId" json:"userId"`
}

func NewItem(title string,description string) Item {
	return Item{
		Title:       title,
		Description: description,
	}
}
