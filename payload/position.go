package payload

import "go.mongodb.org/mongo-driver/bson/primitive"

type Position struct {
	Id primitive.ObjectID  					`bson:"id,omitempty" json:"id"`
	OrganisationId primitive.ObjectID 		`bson:"organisationId,omitempty" json:"organisation_id"`
	Name string 							`bson:"name" json:"name"`
	AccessLevel int 						`bson:"accessLevel" json:"access_level"`
	Token string 							`json:"token"`
}
