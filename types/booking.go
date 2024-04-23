package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	RoomId         primitive.ObjectID `bson:"roomId,omitempty" json:"roomId,omitempty"`
	FromDate       time.Time          `bson:"fromDate,omitempty" json:"fromDate,omitempty"`
	PersonQuantity int                `bson:"personQuantity,omitempty" json:"personQuantity,omitempty"`
	TillDate       time.Time          `bson:"tillDate,omitempty" json:"tillDate,omitempty"`
	UserId         primitive.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"`
	Canceled       bool               `bson:"canceled,omitempty" json:"canceled,omitempty"`
}
