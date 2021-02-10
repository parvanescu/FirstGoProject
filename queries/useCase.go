package queries

import (
	"ExGabi/payload"
	"ExGabi/response"
)

type IUseCase interface {
	GetItemById(item *payload.Item)(*response.Item,error)
	GetItemByTitle(item *payload.Item)(*response.Item,error)
	GetItemByDescription(item *payload.Item)(*[]response.Item,string,error)
	GetAllItems(token string) (*[]response.Item,string,error)

	GetUserById(user *payload.User)(*response.User,error)
	GetUserByEmail(user *payload.User)(*response.User,error)
	GetAllUsers(token string) (*[]response.User,string,error)
	
	Login(user *payload.User)(string,error)
}
