package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/bzawada1/hotel-reservation-app/db"
	"github.com/bzawada1/hotel-reservation-app/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookRoomParams struct {
	FromDate       time.Time `json:"fromDate"`
	TillDate       time.Time `json:"tillDate"`
	PersonQuantity int       `json:"personQuantity`
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleBookHandler(c *fiber.Ctx) error {
	roomOid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	params := BookRoomParams{}
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type:    "error",
			Message: "internal server error",
		})
	}
	booking := types.Booking{
		UserId:         user.ID,
		RoomId:         roomOid,
		FromDate:       params.FromDate,
		TillDate:       params.TillDate,
		PersonQuantity: params.PersonQuantity,
	}

	return c.JSON()
}

func (h *RoomHandler) HandleGetHotel(c *fiber.Ctx) error {
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

func (h *RoomHandler) HandleDeleteHotel(c *fiber.Ctx) error {
	hotelId := c.Params("id")
	if err := h.store.Hotel.DeleteHotel(c.Context(), hotelId); err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": hotelId})
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	hotelId := c.Params("id")
	rooms, err := h.store.Room.GetRoomsByHotelId(c.Context(), hotelId)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}
