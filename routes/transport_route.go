package routes

import (
    "go-trans/controllers" 
    "github.com/gofiber/fiber/v2"
)

func TransportRoute(app *fiber.App) {
	app.Post("/transport", controllers.CreateTransport)
	app.Get("/transport/:transportId", controllers.GetTransport)
	app.Put("/transport/:transportId", controllers.EditTransport)
	app.Delete("/transport/:transportId", controllers.DeleteTransport)
	app.Get("/transports", controllers.GetAllTransports)

}