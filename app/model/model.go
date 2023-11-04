package model

type ObjectID string

type BaseModel struct {
	ID        ObjectID `json:"id" bson:"_id"`
	CreatedAt int64    `json:"createdAt" bson:"createdAt"`
	UpdatedAt int64    `json:"updatedAt" bson:"updatedAt"`
}
