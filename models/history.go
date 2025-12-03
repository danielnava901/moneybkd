package models

import "time"

type History struct {
	ID          int64     `json:"id,omitempty" bson:"_id,omitempty"`
	CountryName string    `json:"country_name" bson:"country_name"`
	CountryCode string    `json:"country_code" bson:"country_code"`
	Value       float64   `json:"value" bson:"value"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}
