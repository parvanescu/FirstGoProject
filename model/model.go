package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ToDoItem struct{
	ItemId      primitive.ObjectID				`bson:"_id,omitempty"`
	Title       string							`bson:"title,omitempty"`
	Description string							`bson:"description,omitempty"`
}

func New(title string,description string) ToDoItem {
	return ToDoItem{
		Title:       title,
		Description: description,
	}
}

func (item ToDoItem)GetItemId()primitive.ObjectID{
	return item.ItemId
}

func (item ToDoItem)GetTitle()string{
	return item.Title
}

func (item ToDoItem)GetDescription()string{
	return item.Description
}


type User struct{
	UserId primitive.ObjectID `bson:"_id,omitempty"`
	UserName string`bson:"name,omitempty"`
	PostedItemIds []primitive.ObjectID `bson:"itemsIds,omitempty"`
}

func (item User)GetUserId()primitive.ObjectID{
	return item.UserId
}

func(item User)GetUserName()string{
	return item.UserName
}

func(item User)GetPostedItemsIds()[]primitive.ObjectID{
	return item.PostedItemIds
}
