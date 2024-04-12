package main

import (
	"flag"

	"github.com/bzawada1/hotel-reservation-app/api"
	"github.com/gofiber/fiber/v2"
)

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()
	app := fiber.New()
	appiv1 := app.Group("/api/v1")

	app.Get("/foo", handleFoo)
	appiv1.Get("/user", api.HandleGetUsers)
	appiv1.Get("/user/:id", api.HandleGetUser)
	app.Listen(*listenAddr)
}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "working just"})
}
