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
	GetItemByTitle(item *payload.Item)(*response.Item,error)
	GetItemByDescription(item *payload.Item)(*[]response.Item,string,error)
	GetAllItems(token string) (*[]response.Item,string,error)

	DeleteUser(user *payload.User)(string,error)
	UpdateUser(user *payload.User)(*response.User,error)
	GetUserById(user *payload.User)(*response.User,error)
	GetUserByEmail(user *payload.User)(*response.User,error)
	GetAllUsers(token string) (*[]response.User,string,error)

	Register(user *payload.User)(string,error)
	Login(user *payload.User)(string,error)
}