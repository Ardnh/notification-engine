package controllers

import (
	"encoding/json"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type WebPushInput struct {
	Headings string `json:"headings" form:"headings" validate:"required,string"`
	Contents string `json:"contents" form:"contents" validate:"required,string"`
}

type TypeInput struct {
	En string `json:"en"`
}

type PostNotificationData struct {
	APP_ID   string    `json:"app_id" validate:"required, string"`
	Headings TypeInput `json:"headings" validate:"required"`
	Contents TypeInput `json:"contents" validate:"required"`
	Segments []string  `json:"included_segments" validate:"required"`
}

func readWebPushEnvKey(key string) string {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error ketika membaca file .env %s", err)
	}
	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("key yang di input salah!")
	}
	return value
}

func webPushUrlBuilder() string {
	webPushAPIBaseURL := readWebPushEnvKey("WEB_PUSH_BASE_API")

	return webPushAPIBaseURL
}

func WebPushHandler(c *fiber.Ctx) error {

	client := resty.New()
	webPushAPIBaseURL := webPushUrlBuilder()
	app_id := readWebPushEnvKey("PUSH_APP_ID")
	api_token := readWebPushEnvKey("API_TOKEN")

	// Parse incoming request
	var input WebPushInput
	err := c.BodyParser(&input)
	headings := &input.Headings
	content := &input.Contents

	// check error during parsing incoming request
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "false",
			"message": err,
		})
	}

	postNotificationData := PostNotificationData{
		APP_ID:   app_id,
		Headings: TypeInput{En: *headings},
		Contents: TypeInput{En: *content},
		Segments: []string{"Subscribed Users"},
	}

	output, _ := json.Marshal(postNotificationData)
	jsonData := string(output)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(api_token).
		SetBody(jsonData).
		Post(webPushAPIBaseURL)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "gagal post notification",
			"message": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  "berhasil",
		"message": resp,
	})
}
