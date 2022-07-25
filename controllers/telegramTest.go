package controllers

import (
	"fmt"
	"log"

	"github.com/goccy/go-json"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type Input struct {
	Sender  string   `json:"sender" form:"sender" validate:"required,string"`
	ChatId  []string `json:"chat_id" form:"chat_id" validate:"required"`
	Text    string   `json:"text,omitempty" form:"message" validate:"string"`
	Caption string   `json:"caption,omitempty" form:"caption" validate:"string"`
}

type PostData struct {
	ChatId string `json:"chat_id"`
	Text   string `json:"text"`
}

func readTelegramEnvKey(key string) string {
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

func urlBuilder(telegramApi string) string {
	const telegramAPIBaseURL string = "https://api.telegram.org/bot"
	bot_token := readTelegramEnvKey("BOT_TOKEN")
	var telegramTokenEnv string = bot_token
	var telegramUrl = telegramAPIBaseURL + telegramTokenEnv + telegramApi
	return telegramUrl
}

func TelegramHandler(c *fiber.Ctx) error {
	client := resty.New()
	const telegramAPISendMessage string = "/sendMessage"
	// parse incoming request
	var input Input
	err := c.BodyParser(&input)
	// check error during parsing incoming request
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "false",
			"message": err,
		})
	}
	// store incoming request data in variable
	var telegramSender = &input.Sender
	var telegramChatId = &input.ChatId
	var telegramMessage = &input.Text
	// var telegramCaption = &input.Caption
	var response []interface{}
	if *telegramMessage != "" {
		for _, value := range *telegramChatId {
			postData := PostData{
				ChatId: value,
				Text:   *telegramMessage,
			}
			output, err := json.Marshal(postData)
			if err != nil {
				fmt.Println(err)
			}
			messageData := string(output)
			telegramUrlForSendMessage := urlBuilder(telegramAPISendMessage)
			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(messageData).
				Post(telegramUrlForSendMessage)

			if err != nil {
				fmt.Println("ini error")
				fmt.Println(err)
			}
			response = append(response, string(resp.Status()))
		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"pengirim": telegramSender,
		"data":     response,
	})
}

// func sendText(chatId []string, message string) interface{} {

// 	client := resty.New()
// 	const telegramAPISendMessage string = "/sendMessage"

// 	var response []interface{}
// 	for _, value := range chatId {
// 		postData := PostData{
// 			ChatId: value,
// 			Text:   message,
// 		}
// 		output, _ := json.Marshal(postData)
// 		messageData := string(output)
// 		telegramUrlForSendMessage := urlBuilder(telegramAPISendMessage)
// 		resp, _ := client.R().
// 			SetHeader("Content-Type", "application/json").
// 			SetBody(messageData).
// 			Post(telegramUrlForSendMessage)
// 		response = append(response, string(resp.Status()))
// 	}

// 	return response
// }
