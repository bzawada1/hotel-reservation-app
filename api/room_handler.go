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
	roomOid, err := primitive.ObjectIDFromHex(roomId)
	if err != nil {
		return err
	}
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
	booking := &types.Booking{
		UserId:         user.ID,
		RoomId:         roomOid,
		FromDate:       params.FromDate,
		TillDate:       params.TillDate,
		PersonQuantity: params.PersonQuantity,
	}

	ok, err = h.isRoomAvailableForBooking(c, roomOid, params)
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Type:    "error",
			Message: fmt.Sprintf("room %s already booked", roomId),
		})
	}

	fmt.Println(c.Context(), booking)
	inserted, err := h.store.Booking.Insert(c.Context(), booking)
	if err != nil {
		return err
	}
	return c.JSON(inserted)
}

func (h *RoomHandler) isRoomAvailableForBooking(c *fiber.Ctx, roomId primitive.ObjectID, params BookRoomParams) (bool, error) {
	bookings, err := h.store.Booking.GetBookings(c.Context(), params.FromDate, params.TillDate, roomId)
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
