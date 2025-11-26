package models

import "time"

type Country struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Symbol    string    `json:"symbol" bson:"symbol"`
	Code      string    `json:"code" bson:"code"`
	Value     float64   `json:"value" bson:"value"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
