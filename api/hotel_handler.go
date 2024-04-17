package api

import (
	"context"
	"errors"

	"github.com/bzawada1/hotel-reservation-app/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

type HotelQueryParams struct {
	Rooms bool
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	qparams := HotelQueryParams{}
	if err := c.QueryParser(&qparams); err != nil {
		return err
	}
	hotels, err := h.store.Hotel.GetHotels(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := context.Background()
	user, err := h.store.Hotel.GetHotelById(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "hotel not found"})
		}
		return err
	}
	return c.JSON(user)
}

func (h *HotelHandler) HandleDeleteHotel(c *fiber.Ctx) error {
	hotelId := c.Params("id")
	if err := h.store.Hotel.DeleteHotel(c.Context(), hotelId); err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": hotelId})
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	hotelId := c.Params("id")
	rooms, err := h.store.Room.GetRoomsByHotelId(c.Context(), hotelId)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}
