package model

type Query struct {
	Limit int64 `json:"limit,omitempty" bson:"limit,omitempty"`
	Skip  int64 `json:"skip,omitempty" bson:"skip,omitempty"`
}
