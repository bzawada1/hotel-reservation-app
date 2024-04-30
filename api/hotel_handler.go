package api

import (
	"context"
	"errors"
	"fmt"

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

type ResourceResponse struct {
	Results int `json:"results"`
	Data    any `json:"data"`
	Page    int `json:"page"`
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	pagination := &db.Pagination{}
	if err := c.QueryParser(pagination); err != nil {
		fmt.Println(err)
		return ErrorBadRequest()
	}
	qparams := HotelQueryParams{}
	if err := c.QueryParser(&qparams); err != nil {
		return err
	}
	hotels, err := h.store.Hotel.GetHotels(c.Context(), pagination)
	if err != nil {
		return err
	}

	return c.JSON(ResourceResponse{Results: len(hotels), Data: hotels, Page: int(pagination.Page)})
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := context.Background()
	user, err := h.store.Hotel.GetHotelById(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrorNotFound("hotel")
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
