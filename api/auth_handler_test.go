package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/bzawada1/hotel-reservation-app/db/fixtures"
	"github.com/gofiber/fiber/v2"
)

func TestAuthenticateSuccess(t *testing.T) {

	tdb := setup(t)
	defer tdb.teardown(t)
	insertedUser := fixtures.AddUser(tdb.store, "John Test", "Dutton", "test@yellowstone.mn", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.store.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    insertedUser.Email,
		Password: fmt.Sprintf("%s_%s", insertedUser.FirstName, insertedUser.LastName),
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {

		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected http status of 200 but good %d", resp.StatusCode)
	}
	response := AuthResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Error(err)
	}
	if response.Token == "" {
		t.Fatalf("expected the JWT token to be present in the auth response")
	}
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, response.User) {
		t.Fatal("expected the user to be the insterted user")
	}
}

func TestAuthenticateWithWrongPasswordFailure(t *testing.T) {

	tdb := setup(t)
	defer tdb.teardown(t)
	insertedUser := fixtures.AddUser(tdb.store, "John Test", "Dutton", "test@yellowstone.mn", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.store.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    insertedUser.Email,
		Password: "wrong_password",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {

		t.Error(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected http status of 400 but good %d", resp.StatusCode)
	}
	genResp := genericResp{}
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}
	if genResp.Type != "error" {
		t.Fatalf("expected generic response type to be error but got %s", genResp.Type)
	}
	if genResp.Message != "invalid credentials" {
		t.Fatalf("expected generic response message to be invalid credentials but got %s", genResp.Type)
	}
}
