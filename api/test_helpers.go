package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"

	"github.com/bzawada1/hotel-reservation-app/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testDb struct {
	client *mongo.Client
	store  *db.Store
}

func setup(t *testing.T) *testDb {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err)
	}
	var (
		DBNAME = os.Getenv("TEST_DBNAME")
		DBURI  = os.Getenv("DB_URI")
	)
	fmt.Println("123", DBNAME, os.Getenv("TEST_DBNAME"))
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(DBURI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client, DBNAME)
	store := &db.Store{
		User:    db.NewMongoUserStore(client, DBNAME),
		Hotel:   hotelStore,
		Room:    db.NewMongoRoomStore(client, hotelStore, DBNAME),
		Booking: db.NewMongoBookingStore(client, DBNAME),
	}

	return &testDb{
		client: client,
		store:  store,
	}
}

func (tdb *testDb) teardown(t *testing.T) {
	if err := tdb.client.Database(os.Getenv("TEST_DBNAME")).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}
