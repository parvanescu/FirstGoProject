package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Organisation struct{
	Id      primitive.ObjectID					`bson:"_id,omitempty" json:"organisation_id"`
	Name       string							`bson:"name" json:"name"`
	CUI string									`bson:"cui" json:"cui"`
	Status bool									`bson:"status" json:"status"`
}