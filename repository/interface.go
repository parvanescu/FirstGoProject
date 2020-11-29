package repository

import (
	"ExGabi/payload"
	"ExGabi/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IRepository interface{
	AddItem(item payload.Item)(primitive.ObjectID,error)
	DeleteItem(userId primitive.ObjectID,id primitive.ObjectID)error
	UpdateItem(id primitive.ObjectID,newItem payload.Item)(response.Item,error)
	GetAllItems() (*[]response.Item,error)
	GetItemById(id primitive.ObjectID) (*response.Item,error)

	AddUser(user payload.User)(primitive.ObjectID,error)
	DeleteUser(id primitive.ObjectID)error
	UpdateUser(id primitive.ObjectID,user payload.User)(response.User,error)
	GetAllUsers() (*[]response.User,error)
	GetUserById(id primitive.ObjectID) (response.User,error)
}
