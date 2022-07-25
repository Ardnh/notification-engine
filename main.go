package main

import (
	"log"
	"notification-engine/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	router.SetupRouter(app)

	log.Fatal(app.Listen(":8080"))
}
