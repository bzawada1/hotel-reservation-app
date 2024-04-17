package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bzawada1/hotel-reservation-app/db"
	"github.com/bzawada1/hotel-reservation-app/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	fmt.Println("seeding the database")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DbUri))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DbName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		User:  db.NewMongoUserStore(client),
		Hotel: db.NewMongoHotelStore(client),
		Room:  db.NewMongoRoomStore(client, hotelStore),
	}
	seedHotel(ctx, store)
	seedUser(ctx, store)
}

func seedHotel(ctx context.Context, store *db.Store) {
	hotel := types.Hotel{
		Name:     "Bellacia",
		Location: "France",
		Rooms:    []primitive.ObjectID{},
		Rating:   5,
	}
	rooms := []types.Room{
		{
			Size:      "small",
			Seaside:   false,
			BasePrice: 123.5,
		},
		{
			Size:      "medium",
			Seaside:   false,
			BasePrice: 123.5,
		},
	}

	insertedHotel, err := store.Hotel.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := store.Room.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertedHotel, insertedRoom)
	}
}

func seedUser(ctx context.Context, store *db.Store) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "John",
		LastName:  "Dutton",
		Email:     "test@yellowstone.mn",
		Password:  "password_montana",
	})

	if err != nil {
		log.Fatal(err)
	}

	if _, err := store.User.CreateUser(ctx, user); err != nil {
		log.Fatal(err)
	}
}
