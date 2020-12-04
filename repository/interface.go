package repository

import (
	"ExGabi/payload"
	"ExGabi/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IRepository interface{
	AddItem(item *payload.Item)(primitive.ObjectID,error)
	DeleteItem(item *payload.Item)error
	UpdateItem(item *payload.Item)(*response.Item,error)
	GetAllItems() (*[]response.Item,error)
	GetItemById(id primitive.ObjectID) (*response.Item,error)

	AddUser(user *payload.User)(primitive.ObjectID,error)
	DeleteUser(user *payload.User)error
	UpdateUser(user *payload.User)(*response.User,error)
	GetAllUsers() (*[]response.User,error)
	GetUserById(id primitive.ObjectID) (*response.User,error)
}
