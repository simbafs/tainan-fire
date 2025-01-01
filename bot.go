package main

import (
	"fmt"
	"log"
	"time"

	"tainanfire/bucket"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type Msg struct {
	event *Event
	msg   *gotgbot.Message
}

type Bot struct {
	bot    *gotgbot.Bot
	chat   int64
	bucket *bucket.Bucket[Msg]
}

func NewBot(apikey string, chat int64) *Bot {
	bot, err := gotgbot.NewBot(apikey, nil)
	if err != nil {
		panic(err)
	}
	return &Bot{
		bot:    bot,
		chat:   chat,
		bucket: bucket.New(60*time.Minute, func(a, b Msg) bool { return a.event.Equal(b.event) }),
	}
}

func (b *Bot) SendEvent(e *Event) error {
	oldMsg, ok := b.bucket.Get(e.ID)
	var text string
	var msg *gotgbot.Message
	var err error

	if !ok {
		// New event
		text = "`新事件\n" + e.String() + "`"
		msg, err = b.bot.SendMessage(b.chat, text, &gotgbot.SendMessageOpts{
			ParseMode: "markdown",
		})
	} else if !oldMsg.event.Equal(e) {
		// update old event
		text = "`事件更新\n" + oldMsg.event.Diff(e) + "\n" + e.String() + "`"
		msg, err = oldMsg.msg.Reply(b.bot, text, &gotgbot.SendMessageOpts{
			ParseMode: "markdown",
		})
	} else {
		return nil
	}

	log.Println(text)

	if err != nil {
		return err
	}

	b.bucket.Set(e.ID, Msg{
		event: e,
		msg:   msg,
	})

	return nil
}

func (b *Bot) GC() {
	l := b.bucket.Len()
	b.bucket.GC()
	if l != b.bucket.Len() {
		log.Println("GC", l, b.bucket.Len())
		b.bot.SendMessage(b.chat, fmt.Sprintf("--debug--\ngc: %d -> %d", l, b.bucket.Len()), &gotgbot.SendMessageOpts{})
	}
}
