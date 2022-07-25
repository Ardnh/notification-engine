package router

import (
	"notification-engine/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App) {
	api := app.Group("/send")

	api.Post("/telegram-send", controllers.TelegramHandler)
	api.Post("/email", controllers.EmailHandler)
	api.Post("/webpush", controllers.WebPushHandler)
}
