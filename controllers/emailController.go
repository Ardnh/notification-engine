package controllers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mailgun/mailgun-go/v3"
	"github.com/spf13/viper"
)

type SendEmail struct {
	Sender    string `json:"sender" validate:"required, string"`
	Subject   string `json:"subject" validate:"required, string"`
	Body      string `json:"body" validate:"required, string"`
	Recipient string `json:"recipient" validate:"required, string"`
}

func readApiKey(key string) string {
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

func EmailHandler(c *fiber.Ctx) error {
	var sendEmail SendEmail = SendEmail{}
	domain := readApiKey("MAIL_DOMAIN")
	apiKey := readApiKey("API_KEY")

	mg := mailgun.NewMailgun(domain, apiKey)

	emailSender := &sendEmail.Sender
	emailSubject := &sendEmail.Subject
	emailBody := &sendEmail.Body
	emailRecipient := &sendEmail.Recipient

	err := c.BodyParser(&sendEmail)

	fmt.Println(&sendEmail)

	if err != nil {
		fmt.Println(err)
	}

	message := mg.NewMessage(*emailSender, *emailSubject, *emailBody, *emailRecipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"id":      id,
			"message": err,
		})
	}

	return c.Status(200).JSON(&fiber.Map{
		"id":     id,
		"result": resp,
	})
}
