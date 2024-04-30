package db

import (
	"os"
)

const (
// DbName     = os.Getenv("DBNAME")
// TestDbName = "hotel-reservation-test"

// DbUri      = "mongodb://localhost:27017"
)

var DbName = os.Getenv("DBNAME")
var DbUri = os.Getenv("DB_URI")

type Pagination struct {
	Limit int64
	Page  int64
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}
