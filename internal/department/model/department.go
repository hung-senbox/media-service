package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Department struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	LocationID         string             `bson:"location_id"`
	OrganizationID     string             `bson:"organization_id"`
	RegionID           string             `bson:"region_id"`
	Icon               string             `bson:"icon"`
	Name               string             `bson:"name"`
	Description        string             `bson:"description"`
	Note               string             `bson:"note"`
	Url                string             `bson:"url"`
	IsPublishedMessage bool               `bson:"is_published_message"`
	Leader             Leader             `bson:"leader"`
	Staffs             []Staff            `bson:"staffs"`
	CreatedAt          time.Time          `bson:"created_at"`
	UpdatedAt          time.Time          `bson:"updated_at"`
}

type Leader struct {
	OwnerID   string `bson:"owner_id"`
	OwnerRole string `bson:"owner_role"`
}

type Staff struct {
	OwnerID   string `bson:"owner_id"`
	OwnerRole string `bson:"owner_role"`
	Index     int    `bson:"index"`
}
