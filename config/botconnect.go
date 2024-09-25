package config

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BotConnect() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {

	bot, err := tgbotapi.NewBotAPI("7512358791:AAFgAI1NIlDeKAJ-k2aYiTto5LQp5ILeRM8")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Get the updates channel
	updates := bot.GetUpdatesChan(u)

	return bot, updates
}
