package main

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
)

type Bot struct {
	bot  *gotgbot.Bot
	chat int64
}

func NewBot(apikey string, chat int64) *Bot {
	bot, err := gotgbot.NewBot(apikey, nil)
	if err != nil {
		panic(err)
	}
	return &Bot{
		bot:  bot,
		chat: chat,
	}
}

func (b *Bot) Send(msg string) (*gotgbot.Message, error) {
	return b.bot.SendMessage(b.chat, msg, nil)
}

func (b *Bot) Reply(msg string, prev *gotgbot.Message) (*gotgbot.Message, error) {
	return prev.Reply(b.bot, msg, nil)
}
