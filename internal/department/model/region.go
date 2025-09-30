package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Region struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	OrganizationID string             `bson:"organization_id"`
	Name           string             `bson:"name"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
}
