package useCase

import (
	"ExGabi/payload"
	"ExGabi/response"
)

type IUseCase interface {
	AddItem(item *payload.Item)(*response.Item,error)
	DeleteItem(item *payload.Item)(string,error)
	UpdateItem(item *payload.Item)(*response.Item,error)
	GetItemById(item *payload.Item)(*response.Item,error)
	GetAllItems() (*[]response.Item,error)

	DeleteUser(user *payload.User)(string,error)
	UpdateUser(user *payload.User)(*response.User,error)
	GetUserById(user *payload.User)(*response.User,error)
	GetAllUsers() (*[]response.User,error)

	Register(user *payload.User)(string,error)
	Login(user *payload.User)(string,error)
}