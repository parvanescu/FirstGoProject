package repository

import (
	"ExGabi/payload"
	"ExGabi/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IRepository interface{
	AddItem(item payload.Item)primitive.ObjectID
	DeleteItem(userId primitive.ObjectID,id primitive.ObjectID)error
	UpdateItem(id primitive.ObjectID,newItem payload.Item)(response.Item,error)
	GetAllItems() *[]response.Item
	GetItemById(id primitive.ObjectID) *response.Item

	AddUser(user payload.User)primitive.ObjectID
	DeleteUser(id primitive.ObjectID)error
	UpdateUser(id primitive.ObjectID,user payload.User)(response.User,error)
	GetAllUsers() *[]response.User
	GetUserById(id primitive.ObjectID) response.User

	GetUserStatus(id primitive.ObjectID)(bool,error)
	SetUserStatus(id primitive.ObjectID,status bool)error
}
