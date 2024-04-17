package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bzawada1/hotel-reservation-app/db"
	"github.com/bzawada1/hotel-reservation-app/types"
	"github.com/gofiber/fiber/v2"
)

func TestAuthenticate(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	insertTestUser(t, tdb.Store)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "james@foo.com",
		Password: "password123",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("expected http status of 200 but good %d", http.StatusOK)
	}
	response := AuthResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&resp); err != nil {
		t.Error(err)
	}

}

func insertTestUser(t *testing.T, store *db.Store) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "John Test",
		LastName:  "Dutton",
		Email:     "test@yellowstone.mn",
		Password:  "password_montana",
	})

	if err != nil {
		t.Fatal(err)
	}

	if _, err := store.User.CreateUser(context.TODO(), user); err != nil {
		t.Fatal(err)
	}
}
