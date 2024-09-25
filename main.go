package main

import (
	"core-bot/config"
	"core-bot/repository/sqlcontroller"
)

func main() {
	config.ConnectDB()
	bot, updates := config.BotConnect()
	var userState = make(map[int64]string)
	sqlcontroller.SqlCommands(bot, updates, userState)

}
