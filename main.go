package main

import (
	"go-trans/configs"
	"go-trans/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "Welcome to transportation"})
	})

	//run database
	configs.ConnectDB()

	//routes
	routes.TransportRoute(app)

	//listen on port
	app.Listen(":6000")
}
