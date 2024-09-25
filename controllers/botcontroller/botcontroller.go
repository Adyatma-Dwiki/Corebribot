package botcontroller

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendGreeting(bot *tgbotapi.BotAPI, update tgbotapi.Update, state map[int64]string) {
	userName := update.Message.From.UserName
	if userName == "" {
		userName = update.Message.From.FirstName
	}

	greeting := fmt.Sprintf("Halo %s, menu apa yang ingin Anda gunakan?\n 1. Query DDMAST \n 2. Query CDMAST \n 3. Query LNMAST", userName)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, greeting)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Failed to send greeting:", err)
		return
	}
	state[update.Message.Chat.ID] = "halo"
}

func SendQueryResult(bot *tgbotapi.BotAPI, update tgbotapi.Update, resultMessage string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, resultMessage)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Failed to send query result:", err)
	}
}
