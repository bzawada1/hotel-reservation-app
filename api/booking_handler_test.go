package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bzawada1/hotel-reservation-app/db/fixtures"
	"github.com/bzawada1/hotel-reservation-app/types"
)

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := fixtures.AddUser(tdb.store, "Alex", "Spencer", "alex-spencer@hotmail.com", false)
	adminUser := fixtures.AddUser(tdb.store, "Alex", "Spencer", "alex-spencer@hotmail.com", true)
	hotel := fixtures.AddHotel(tdb.store, "Paris", "Hilton", 4, nil)
	room := fixtures.AddRoom(tdb.store, "Double bed", true, 140.5, hotel.ID)
	booking := fixtures.AddBooking(tdb.store, room.ID, user.ID)

	app := NewTestApp()
	admin := app.Group("/", JWTAuthentication(tdb.store.User), AdminAuth)
	BookingHandler := NewBookingHandler(tdb.store)
	admin.Get("/", BookingHandler.HandleGetBookings)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Content-Type", "text/json")
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response, received: %d", resp.StatusCode)
	}
	bookings := []*types.Booking{}
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking got %d", len(bookings))
	}
	if booking.ID != bookings[0].ID {
		t.Fatalf("expected booking with id: %d received %d", booking.ID, bookings[0].ID)
	}

	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Content-Type", "text/json")
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status unauthorized, received: %d", resp.StatusCode)
	}
}

func TestUserGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := fixtures.AddUser(tdb.store, "Alex", "Spencer", "alex-spencer@hotmail.com", false)
	hotel := fixtures.AddHotel(tdb.store, "Paris", "Hilton", 4, nil)
	room := fixtures.AddRoom(tdb.store, "Double bed", true, 140.5, hotel.ID)
	booking := fixtures.AddBooking(tdb.store, room.ID, user.ID)

	app := NewTestApp()
	authGroup := app.Group("/", JWTAuthentication(tdb.store.User))
	BookingHandler := NewBookingHandler(tdb.store)
	authGroup.Get("/:id", BookingHandler.HandleGetBooking)

	req := httptest.NewRequest("GET", "/"+booking.ID.Hex(), nil)
	req.Header.Add("Content-Type", "text/json")
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {

		t.Fatalf("non 200 response, received: %d", resp.StatusCode)
	}
	receivedBooking := types.Booking{}
	if err := json.NewDecoder(resp.Body).Decode(&receivedBooking); err != nil {
		t.Fatal(err)
	}
	if booking.ID != receivedBooking.ID {
		t.Fatalf("expected booking with id: %d received %d", booking.ID, receivedBooking.ID)
	}
	if booking.UserId != receivedBooking.UserId {
		t.Fatalf("expected user with id: %d received %d", booking.UserId, receivedBooking.UserId)
	}
}
