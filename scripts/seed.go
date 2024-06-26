package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/bzawada1/hotel-reservation-app/api"
	"github.com/bzawada1/hotel-reservation-app/db"
	"github.com/bzawada1/hotel-reservation-app/db/fixtures"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	fmt.Println("seeding the database")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DbUri))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DbName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client, db.DbName)
	store := &db.Store{
		User:    db.NewMongoUserStore(client, db.DbName),
		Hotel:   hotelStore,
		Room:    db.NewMongoRoomStore(client, hotelStore, db.DbName),
		Booking: db.NewMongoBookingStore(client, db.DbName),
	}

	user := fixtures.AddUser(store, "John", "Dutton", "test2@yellowstone.mn", false)
	admin := fixtures.AddUser(store, "John", "Dutton Admin", "test@yellowstone.mn", true)
	fmt.Println("admin token -->", api.CreateTokenFromUser(admin))
	fmt.Println("user token -->", api.CreateTokenFromUser(user))

	for i := 0; i < 100; i++ {
		hotel := fixtures.AddHotel(store, "Paris", "Hilton", rand.Intn(5)+1, nil)
		room := fixtures.AddRoom(store, "Double bed", true, 140.5, hotel.ID)
		fixtures.AddBooking(store, room.ID, user.ID)
	}
}
