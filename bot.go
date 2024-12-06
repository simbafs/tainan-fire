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
	return b.bot.SendMessage(b.chat, msg, &gotgbot.SendMessageOpts{
		ParseMode: "markdown",
	})
}

func (b *Bot) Reply(msg string, prev *gotgbot.Message) (*gotgbot.Message, error) {
	return prev.Reply(b.bot, msg, nil)
}

func (b *Bot) SendEvent(prev *gotgbot.Message, e *Event) error {
	s := ""
	isUpdate := false
	if prev == nil {
		s += "新事件\n"
	} else {
		isUpdate = true
		s += "事件更新\n"
	}

	s += e.String()

	var m *gotgbot.Message
	var err error
	if isUpdate {
		m, err = b.Reply(s, prev)
	} else {
		m, err = b.Send(s)
	}

	if err != nil {
		return err
	}

	e.PrevMessage = m

	return nil
}
