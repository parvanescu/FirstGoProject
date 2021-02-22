package mutations

import (
	"ExGabi/payload"
	"ExGabi/response"
)

type IUseCase interface {
	AddItem(item *payload.Item)(*response.Item,error)
	DeleteItem(item *payload.Item)(string,error)
	UpdateItem(item *payload.Item)(*response.Item,error)

	DeleteUser(user *payload.User)(string,error)
	UpdateUser(user *payload.User)(*response.User,error)

	Register(user *payload.User)(*response.User,error)
	Login(user *payload.User)(string,error)

	GetMatchingSearch(item *payload.Item)(*[]response.Item,string,error)
}
