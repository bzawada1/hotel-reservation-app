package db

import (
	"context"
	"fmt"

	"github.com/bzawada1/hotel-reservation-app/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const hotelColl = "hotels"

type HotelStore interface {
	Insert(context.Context, *types.Hotel) (*types.Hotel, error)
	GetHotels(context.Context) ([]*types.Hotel, error)
	GetHotelById(context.Context, string) (*types.Hotel, error)
	DeleteHotel(context.Context, string) error
	UpdateHotel(context.Context, bson.M, bson.M) error
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client, dbName string) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(dbName).Collection(hotelColl),
	}
}

func (s *MongoHotelStore) Insert(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	result, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}

	hotel.ID = result.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (s *MongoHotelStore) GetHotels(ctx context.Context) ([]*types.Hotel, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel
	if err := cur.All(ctx, &hotels); err != nil {
		return nil, err
	}

	return hotels, nil
}

func (s *MongoHotelStore) DeleteHotel(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	result, err := s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("couldn't delete hotel with id: %v", id)
	}

	return nil
}

func (s *MongoHotelStore) GetHotelById(ctx context.Context, id string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	hotel := types.Hotel{}
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel); err != nil {
		return nil, err
	}
	return &hotel, nil
}

func (s *MongoHotelStore) UpdateHotel(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}
