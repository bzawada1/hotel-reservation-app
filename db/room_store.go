package db

import (
	"context"

	"github.com/bzawada1/hotel-reservation-app/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomColl = "rooms"

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRoomsByHotelId(context.Context, string) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	hotelStore HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore, dbName string) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(dbName).Collection(roomColl),
		hotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	result, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = result.InsertedID.(primitive.ObjectID)

	filter := Filter{"_id": room.HotelID}
	update := Filter{"$push": Filter{"rooms": room.ID}}
	if err := s.hotelStore.UpdateHotel(ctx, filter, update); err != nil {
		return nil, err
	}
	return room, nil
}

func (s *MongoRoomStore) GetRoomsByHotelId(ctx context.Context, hotelId string) ([]*types.Room, error) {
	hotelOid, err := primitive.ObjectIDFromHex(hotelId)
	if err != nil {
		return nil, err
	}
	cur, err := s.coll.Find(ctx, Filter{"hotelID": hotelOid})
	if err != nil {
		return nil, err
	}
	rooms := []*types.Room{}
	if err := cur.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}
