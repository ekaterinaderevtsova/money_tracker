package handler

import (
	"gopkg.in/telebot.v4"
)

func (th *TelegramHandler) SetCommands(bot *telebot.Bot) {
	bot.Handle("/add", th.SpendingHandler.AddSpending)
	bot.Handle("/getweek", th.SpendingHandler.GetWeekSpendings)
	bot.Handle("/getmonth", th.SpendingHandler.GetMonthSpendings)
}
