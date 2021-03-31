package model

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct{
	mgm.DefaultModel 							`bson:",inline"`
	Id      primitive.ObjectID					`bson:"_id,omitempty" json:"item_id"`
	Title       string							`bson:"title" json:"title"`
	Description string							`bson:"description" json:"description"`
	UserId 		primitive.ObjectID				`bson:"userId" json:"userId"`
}


