package model

type Email struct {
	BaseModel `bson:",inline"`

	To      []string `json:"to" bson:"to"`
	From    string   `json:"from" bson:"from"`
	Subject string   `json:"subject" bson:"subject"`

	Headers map[string]string `json:"headers" bson:"headers"`
	Body    string            `json:"body" bson:"body"`
}
