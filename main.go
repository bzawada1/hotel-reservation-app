package main

import (
	"context"
	"flag"
	"log"

	"github.com/bzawada1/hotel-reservation-app/api"
	"github.com/bzawada1/hotel-reservation-app/api/middleware"
	"github.com/bzawada1/hotel-reservation-app/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()
	app := fiber.New(config)
	apiv1 := app.Group("/api/v1", middleware.JWTAuthentication)
	auth := app.Group("/api")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DbUri))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		User:  db.NewMongoUserStore(client),
		Hotel: hotelStore,
		Room:  db.NewMongoRoomStore(client, hotelStore),
	}
	userHandler := api.NewUserHandler(store)
	hotelHandler := api.NewHotelHandler(store)
	authHandler := api.NewAuthHandler(store.User)

	auth.Post("/auth", authHandler.HandleAuthenticate)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	app.Listen(*listenAddr)
}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "working just"})
}
