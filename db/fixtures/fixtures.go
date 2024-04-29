package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bzawada1/hotel-reservation-app/db"
	"github.com/bzawada1/hotel-reservation-app/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.Store, fistName, lastName, email string, isAdmin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fistName,
		LastName:  lastName,
		Email:     email,
		Password:  fmt.Sprintf("%s_%s", fistName, lastName),
	})

	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	insertedUser, err := store.User.CreateUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	return insertedUser
}

func AddHotel(store *db.Store, location, name string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	roomsIds := rooms
	if rooms == nil {
		roomsIds = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    roomsIds,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.Insert(context.Background(), &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func AddRoom(store *db.Store, size string, seaside bool, price float64, hotelId primitive.ObjectID) *types.Room {
	room := types.Room{
		Size:    size,
		Seaside: seaside,
		Price:   price,
		HotelID: hotelId,
	}

	insertedRoom, err := store.Room.InsertRoom(context.Background(), &room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func AddBooking(store *db.Store, roomId primitive.ObjectID, userId primitive.ObjectID) *types.Booking {
	booking := &types.Booking{
		RoomId:         roomId,
		FromDate:       time.Now(),
		PersonQuantity: 4,
		TillDate:       time.Now().AddDate(0, 0, 2),
		UserId:         userId,
	}
	insertedBooking, err := store.Booking.Insert(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}
