package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/bzawada1/hotel-reservation-app/types"
)

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := NewTestApp()
	UserHandler := NewUserHandler(tdb.store)
	app.Post("/", UserHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "some@fo.com",
		FirstName: "James",
		LastName:  "McGuilll",
		Password:  "asfadsfgvs",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "text/json")
	resp, _ := app.Test(req)
	user := types.User{}
	json.NewDecoder(resp.Body).Decode(&user)
	if len(user.ID) == 0 {
		t.Errorf("expecting a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expecting the encrypted password not to be included in the json response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstname %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastname %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}

	fmt.Println(resp.Status)
}
