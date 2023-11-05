package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ObjectID = primitive.ObjectID

type BaseModel struct {
	ID        ObjectID  `json:"id" bson:"_id"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}
