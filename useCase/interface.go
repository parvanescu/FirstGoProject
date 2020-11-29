package useCase

import (
	"ExGabi/payload"
	"ExGabi/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUseCase interface {
	AddItem(item payload.Item)(primitive.ObjectID,error)
	DeleteItem(userId primitive.ObjectID,itemId primitive.ObjectID)error
	UpdateItem(id primitive.ObjectID,item payload.Item)(response.Item,error)
	GetItem(id primitive.ObjectID)(response.Item,error)
	GetAllItems() (*[]response.Item,error)

	AddUser(user payload.User)(primitive.ObjectID,error)
	DeleteUser(id primitive.ObjectID)error
	UpdateUser(id primitive.ObjectID,user payload.User)(response.User,error)
	GetUser(id primitive.ObjectID)(response.User,error)
	GetAllUsers() (*[]response.User,error)
}