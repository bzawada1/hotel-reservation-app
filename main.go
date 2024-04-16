package main

import (
	"context"
	"flag"
	"log"

	"github.com/bzawada1/hotel-reservation-app/api"
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
	apiv1 := app.Group("/api/v1")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DbUri))
	if err != nil {
		log.Fatal(err)
	}
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client, db.DbName))
	var hotelStore = db.NewMongoHotelStore(client)
	var roomStore = db.NewMongoRoomStore(client, hotelStore)
	var hotelHandler = api.NewHotelHandler(hotelStore, roomStore)

	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	app.Listen(*listenAddr)
}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "working just"})
}
