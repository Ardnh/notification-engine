package main

import (
	"notification-engine/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	router.SetupRouter(app)

	app.Listen(":8080")
}
