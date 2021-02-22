package queries

import (
	"ExGabi/payload"
	"ExGabi/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IRepository interface{
	GetAllItems() (*[]response.Item,error)
	GetAllUsersItems(userId primitive.ObjectID) (*[]response.Item,error)
	GetItemByTitle(item *payload.Item)(*response.Item,error)
	GetItemByDescription(item *payload.Item)(*[]response.Item,error)
	GetItemById(id primitive.ObjectID) (*response.Item,error)

	GetAllUsers() (*[]response.User,error)
	GetUserById(id primitive.ObjectID) (*response.User,error)
	GetUserByCredentials(user *payload.User)(*response.User,error)
	GetUserByEmail(user *payload.User)(*response.User,error)
}