package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName string             `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName  string             `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Username  string             `json:"username,omitempty" bson:"username,omitempty,unique"`
	Password  string             `json:"password,omitempty" bson:"password,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type CreateUser struct {
	Username string `json:"username,omitempty" bson:"username,omitempty,unique" validate:"required,min=6,max=30"`
	Password string `json:"password,omitempty" bson:"password,omitempty" validate:"required,min=8,max=30,containsany=!@#$$%^&*(),uppercase,lowercase,number"`
}
