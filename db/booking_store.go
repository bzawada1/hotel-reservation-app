package db

import (
	"context"
	"fmt"
	"time"

	"github.com/bzawada1/hotel-reservation-app/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookingColl = "bookings"

type BookingStore interface {
	Insert(context.Context, *types.Booking) (*types.Booking, error)
	GetAllBookings(context.Context, time.Time, time.Time) ([]*types.Booking, error)
	GetBookings(context.Context, time.Time, time.Time, primitive.ObjectID) ([]*types.Booking, error)
	GetBookingById(context.Context, string) (*types.Booking, error)
	GetRooms(context.Context) ([]*types.Booking, error)
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

func (s *MongoBookingStore) Insert(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	result, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}

	booking.ID = result.InsertedID.(primitive.ObjectID)
	return booking, nil
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, fromDate time.Time, tillDate time.Time, roomId primitive.ObjectID) ([]*types.Booking, error) {

	filter := bson.M{
		"roomId": roomId,
		"fromDate": bson.M{
			"$gte": primitive.NewDateTimeFromTime(fromDate),
		},
		"tillDate": bson.M{
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
	fmt.Println(booking)

	return booking, nil
}

func (s *MongoBookingStore) GetAllBookings(ctx context.Context, fromDate time.Time, tillDate time.Time) ([]*types.Booking, error) {

	filter := bson.M{}
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
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking); err != nil {
		return nil, err
	}
	return &booking, nil
}

func (s *MongoBookingStore) GetRooms(ctx context.Context) ([]*types.Booking, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	bookings := []*types.Booking{}
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}
