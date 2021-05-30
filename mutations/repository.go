package mutations

import (
	"ExGabi/payload"
	"ExGabi/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IRepository interface{
	AddItem(userId primitive.ObjectID,item *payload.Item)(primitive.ObjectID,error)
	DeleteItem(userId primitive.ObjectID,item *payload.Item)error
	UpdateItem(userId primitive.ObjectID,item *payload.Item)(*response.Item,error)

	AddUser(organisationId primitive.ObjectID,user *payload.User)(primitive.ObjectID,error)
	DeleteUser(user *payload.User)error
	UpdateUser(user *payload.User)(*response.User,error)

	AddOrganisation(organisation *payload.Organisation)(primitive.ObjectID,error)
	DeleteOrganisationById(organisation *payload.Organisation)error
	UpdateOrganisation(organisation *payload.Organisation)(*response.Organisation,error)

	GetUserById(id primitive.ObjectID) (*response.User,error)
	GetUserByEmail(user *payload.User)(*response.User,error)
	GetUserByCredentials(user *payload.User)(*response.User,error)
	GetMatchingItems(userId primitive.ObjectID,item *payload.Item)(*[]response.Item,error)
	GetOrganisationByCUI(organisation *payload.Organisation)(*response.Organisation,error)

	AddPositionToOrganisation(position *payload.Position)(primitive.ObjectID,error)
	UpdatePosition(position *payload.Position)(*response.Position,error)
	ExchangePositions(position1 *payload.Position,position2 *payload.Position)(*response.Position,*response.Position,error)

	AddUserAndOrganisation(user *payload.User,organisation *payload.Organisation)(primitive.ObjectID,primitive.ObjectID,error)
	UpdateUserPassword(user *payload.User,organisation *payload.Organisation) error
}
