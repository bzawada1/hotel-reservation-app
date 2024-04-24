package api

import (
	"context"
	"log"
	"testing"

	"github.com/bzawada1/hotel-reservation-app/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testDbUri  = "mongodb://localhost:27017"
	testDbName = "hotel-reservation-test"
)

type testDb struct {
	client *mongo.Client
	store  *db.Store
}

func setup(t *testing.T) *testDb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DbUri))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client, db.TestDbName)
	store := &db.Store{
		User:    db.NewMongoUserStore(client, db.TestDbName),
		Hotel:   hotelStore,
		Room:    db.NewMongoRoomStore(client, hotelStore, db.TestDbName),
		Booking: db.NewMongoBookingStore(client, db.DbName),
	}

	return &testDb{
		client: client,
		store:  store,
	}
}

func (tdb *testDb) teardown(t *testing.T) {
	if err := tdb.client.Database(testDbName).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}
