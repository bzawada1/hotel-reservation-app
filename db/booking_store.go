package db

import (
	"context"
	"time"

	"github.com/bzawada1/hotel-reservation-app/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookingColl = "bookings"

type BookingStore interface {
	Insert(context.Context, *types.Booking, string) (*types.Booking, error)
	GetAllBookings(context.Context, time.Time, time.Time) ([]*types.Booking, error)
	GetBookings(context.Context, time.Time, time.Time, primitive.ObjectID) ([]*types.Booking, error)
	GetBookingById(context.Context, string) (*types.Booking, error)
	GetRooms(context.Context) ([]*types.Booking, error)
	updateBooking(context.Context, primitive.ObjectID, Filter) error
	CancelBooking(context.Context, *types.Booking) error
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client, dbName string) *MongoBookingStore {
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(dbName).Collection(bookingColl),
	}
}

func (s *MongoBookingStore) Insert(ctx context.Context, booking *types.Booking, roomId string) (*types.Booking, error) {
	roomOid, err := primitive.ObjectIDFromHex(roomId)
	if err != nil {
		return nil, err
	}
	booking.RoomId = roomOid
	result, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}

	booking.ID = result.InsertedID.(primitive.ObjectID)
	return booking, nil
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, fromDate time.Time, tillDate time.Time, roomId primitive.ObjectID) ([]*types.Booking, error) {

	filter := Filter{
		"roomId": roomId,
		"fromDate": Filter{
			"$gte": primitive.NewDateTimeFromTime(fromDate),
		},
		"tillDate": Filter{
			"$lte": primitive.NewDateTimeFromTime(tillDate),
		},
	}
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	booking := []*types.Booking{}
	if err := cur.All(ctx, &booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *MongoBookingStore) GetAllBookings(ctx context.Context, fromDate time.Time, tillDate time.Time) ([]*types.Booking, error) {
	filter := Filter{}
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	booking := []*types.Booking{}
	if err := cur.All(ctx, &booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *MongoBookingStore) GetBookingById(ctx context.Context, id string) (*types.Booking, error) {
	booking := types.Booking{}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	if err := s.coll.FindOne(ctx, Filter{"_id": oid}).Decode(&booking); err != nil {
		return nil, err
	}
	return &booking, nil
}

func (s *MongoBookingStore) GetRooms(ctx context.Context) ([]*types.Booking, error) {
	cur, err := s.coll.Find(ctx, Filter{})
	if err != nil {
		return nil, err
	}
	bookings := []*types.Booking{}
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (s *MongoBookingStore) CancelBooking(c context.Context, b *types.Booking) error {
	oid, err := primitive.ObjectIDFromHex(b.ID.Hex())
	if err != nil {
		return err
	}
	updateParams := Filter{"canceled": true}
	s.updateBooking(c, oid, updateParams)
	return err
}

func (s *MongoBookingStore) updateBooking(c context.Context, id primitive.ObjectID, updateParams Filter) error {
	update := bson.D{
		{
			"$set", updateParams,
		},
	}
	_, err := s.coll.UpdateByID(c, id, update)
	return err
}
