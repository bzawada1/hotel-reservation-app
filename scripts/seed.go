package main

import (
	"context"
	"fmt"
	"log"
	"time"

	// "time"

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
	hotelStore := db.NewMongoHotelStore(client, db.DbName)
	store := &db.Store{
		User:    db.NewMongoUserStore(client, db.DbName),
		Hotel:   hotelStore,
		Room:    db.NewMongoRoomStore(client, hotelStore, db.DbName),
		Booking: db.NewMongoBookingStore(client, db.DbName),
	}

	hotelId := seedHotel(ctx, store)
	roomId := seedRoom(ctx, store, hotelId)
	userId := seedUser(ctx, store, "test@yellowstone.mn", false)
	seedBooking(ctx, store, roomId, userId)
	seedUser(ctx, store, "test2@yellowstone.mn", true)
}

func seedHotel(ctx context.Context, store *db.Store) primitive.ObjectID {
	hotel := types.Hotel{
		Name:     "Bellacia",
		Location: "France",
		Rooms:    []primitive.ObjectID{},
		Rating:   5,
	}

	insertedHotel, err := store.Hotel.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel.ID
}

func seedRoom(ctx context.Context, store *db.Store, hotelId primitive.ObjectID) primitive.ObjectID {
	room := types.Room{
		Size:      "small",
		Seaside:   false,
		BasePrice: 123.5,
		HotelID:   hotelId,
	}
	insertedRoom, err := store.Room.InsertRoom(ctx, &room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom.ID
}

func seedBooking(ctx context.Context, store *db.Store, roomId primitive.ObjectID, userId primitive.ObjectID) primitive.ObjectID {
	booking := &types.Booking{
		RoomId:         roomId,
		FromDate:       time.Now(),
		PersonQuantity: 90,
		TillDate:       time.Now().AddDate(0, 0, 2),
		UserId:         userId,
	}
	insertedBooking, err := store.Booking.Insert(ctx, booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking.ID
}

func seedUser(ctx context.Context, store *db.Store, email string, isAdmin bool) primitive.ObjectID {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "John",
		LastName:  "Dutton",
		Email:     email,
		Password:  "password_montana",
	})

	if err != nil {
		log.Fatal(err)
	}
	if isAdmin {
		user.IsAdmin = true
	}
	insertedUser, err := store.User.CreateUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	return insertedUser.ID
}
