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
	roomStore := db.NewMongoRoomStore(client, hotelStore)
	hotel := types.Hotel{
		Name:     "Bellacia",
		Location: "France",
		Rooms:    []primitive.ObjectID{},
		Rating:   5,
	}
	rooms := []types.Room{
		{
			Type:      types.SinglePersonRoomType,
			BasePrice: 123.5,
		},
		{
			Type:      types.DeluxeRoomType,
			BasePrice: 123.5,
		},
	}

	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertedHotel, insertedRoom)
	}
}
