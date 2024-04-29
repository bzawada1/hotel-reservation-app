package api

import (
	"context"
	"errors"

	"github.com/bzawada1/hotel-reservation-app/db"
	"github.com/bzawada1/hotel-reservation-app/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	store *db.Store
}

func NewUserHandler(store *db.Store) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.store.User.GetUsers(c.Context())
	if err != nil {
		return ErrorNotFound("users")
	}
	return c.JSON(users)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := context.Background()
	user, err := h.store.User.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrorNotFound("user")
		}
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUserByEmail(c *fiber.Ctx) error {
	email := c.Params("Email")
	ctx := context.Background()
	user, err := h.store.User.GetUserById(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrorNotFound("user")
		}
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	params := types.CreateUserParams{}
	if err := c.BodyParser(&params); err != nil {
		return ErrorBadRequest()
	}
	if errors := params.Validate(); len(errors) > 0 {
		return ErrorBadRequest()
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return ErrorBadRequest()
	}
	insertedUser, err := h.store.User.CreateUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	params := types.UpdateUserParams{}
	userId := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return ErrorBadRequest()
	}
	if err := c.BodyParser(&params); err != nil {
		return ErrorBadRequest()
	}
	filter := bson.M{"_id": oid}
	if err := h.store.User.UpdateUser(c.Context(), filter, params); err != nil {
		return err
	}
	return c.JSON(map[string]string{"updated": userId})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	if err := h.store.User.DeleteUser(c.Context(), userId); err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": userId})
}
