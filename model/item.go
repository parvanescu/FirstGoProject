package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct{
	ItemId      primitive.ObjectID				`bson:"_id,omitempty"`
	Title       string							`bson:"title,omitempty"`
	Description string							`bson:"description,omitempty"`
}

func NewItem(title string,description string) Item {
	return Item{
		Title:       title,
		Description: description,
	}
}


