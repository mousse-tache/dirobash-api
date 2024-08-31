package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Citation struct {
	Id     primitive.ObjectID `json:"_id,omitempty"`
	Date   primitive.DateTime `json:"date,omitempty" validate:"required"`
	Number int                `json:"number,omitempty" validate:"required"`
	Text   string             `json:"text,omitempty" validate:"required"`
}
