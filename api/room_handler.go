package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bzawada1/hotel-reservation-app/db"
	"github.com/bzawada1/hotel-reservation-app/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate       time.Time `json:"fromDate"`
	TillDate       time.Time `json:"tillDate"`
	PersonQuantity int       `json:"personQuantity"`
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return NewError(500, fmt.Sprintf("cannot book a room in the past"))
	}
	return nil
}

func (h *RoomHandler) HandleBooking(c *fiber.Ctx) error {
	roomId := c.Params("id")
	params := BookRoomParams{}
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.validate(); err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResp{
			Type:    "error",
			Message: "internal server error",
		})
	}

	ok, err := h.isRoomAvailableForBooking(c, roomId, user, params)
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Type:    "error",
			Message: fmt.Sprintf("room %s already booked", roomId),
		})
	}

	booking := &types.Booking{
		UserId:         user.ID,
		FromDate:       params.FromDate,
		TillDate:       params.TillDate,
		PersonQuantity: params.PersonQuantity,
	}
	inserted, err := h.store.Booking.Insert(c.Context(), booking, roomId)
	if err != nil {
		return err
	}
	return c.JSON(inserted)
}

func (h *RoomHandler) isRoomAvailableForBooking(c *fiber.Ctx, roomId string, user *types.User, params BookRoomParams) (bool, error) {
	roomOid, err := primitive.ObjectIDFromHex(roomId)
	if err != nil {
		return false, err
	}
	bookings, err := h.store.Booking.GetBookings(c.Context(), params.FromDate, params.TillDate, roomOid)
	if err != nil {
		return false, err
	}
	ok := len(bookings) == 0
	return ok, nil
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	qparams := BookRoomParams{}
	if err := c.QueryParser(&qparams); err != nil {
		return err
	}
	rooms, err := h.store.Booking.GetRooms(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(rooms)
}
