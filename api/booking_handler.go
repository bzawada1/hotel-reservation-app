package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bzawada1/hotel-reservation-app/db"
	"github.com/bzawada1/hotel-reservation-app/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	fromDate := time.Now().AddDate(0, -5, 0)
	toDate := time.Now()
	fmt.Println(fromDate, toDate)
	bookings, err := h.store.Booking.GetAllBookings(c.Context(), fromDate, toDate)

	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingById(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "booking not found"})
		}
		return err
	}
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return err
	}
	if booking.UserId != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type:    "error",
			Message: "not authorized",
		})
	}
	return c.JSON(booking)
}
