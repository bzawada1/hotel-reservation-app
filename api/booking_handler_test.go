package api

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/bzawada1/hotel-reservation-app/api/middleware"
	"github.com/bzawada1/hotel-reservation-app/db/fixtures"
	"github.com/bzawada1/hotel-reservation-app/types"
	"github.com/gofiber/fiber/v2"
)

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := fixtures.AddUser(tdb.store, "Alex", "Spencer", "alex-spencer@hotmail.com", true)
	hotel := fixtures.AddHotel(tdb.store, "Paris", "Hilton", 4, nil)
	room := fixtures.AddRoom(tdb.store, "Double bed", true, 140.5, hotel.ID)
	fixtures.AddBooking(tdb.store, room.ID, user.ID)

	app := fiber.New()
	admin := app.Group("/admin", middleware.AdminAuth)
	BookingHandler := NewBookingHandler(tdb.store)
	admin.Get("/", BookingHandler.HandleGetBookings)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Content-Type", "text/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	bookings := []*types.Booking{}
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

}
