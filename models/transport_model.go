package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Transport struct {
    Id       	primitive.ObjectID `json:"id,omitempty"`
    Name     	string             `json:"name,omitempty" validate:"required"`
    Description string             `json:"description,omitempty" validate:"required"`
    Modality    string             `json:"modality,omitempty" validate:"required"`
}