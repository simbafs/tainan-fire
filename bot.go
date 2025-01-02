package main

import (
	"fmt"
	"log"
	"strings"
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

type BotOpt struct {
	APIKey    string
	ChatID    int64
	AliveTime time.Duration
}

func WithAPIKey(apikey string) func(*BotOpt) {
	return func(opt *BotOpt) {
		opt.APIKey = apikey
	}
}

func WithChatID(chat int64) func(*BotOpt) {
	return func(opt *BotOpt) {
		opt.ChatID = chat
	}
}

func WithAliveTime(d time.Duration) func(*BotOpt) {
	return func(opt *BotOpt) {
		opt.AliveTime = d
	}
}

func NewBot(opts ...func(*BotOpt)) *Bot {
	defaultOpt := &BotOpt{
		AliveTime: 48 * time.Hour,
	}

	for _, opt := range opts {
		opt(defaultOpt)
	}

	bot, err := gotgbot.NewBot(defaultOpt.APIKey, nil)
	if err != nil {
		panic(err)
	}
	return &Bot{
		bot:    bot,
		chat:   defaultOpt.ChatID,
		bucket: bucket.New(defaultOpt.AliveTime, func(a, b Msg) bool { return a.event.Equal(b.event) }),
	}
}

func (b *Bot) SendMessage(msg string) (*gotgbot.Message, error) {
	// escape markdown
	msg = strings.ReplaceAll(msg, "-", "\\-")
	return b.bot.SendMessage(b.chat, msg, &gotgbot.SendMessageOpts{
		ParseMode: "MarkdownV2",
	})
}

func (b *Bot) SendEvent(e *Event) error {
	oldMsg, ok := b.bucket.Get(e.ID)
	var text string
	var msg *gotgbot.Message
	var err error

	if !ok {
		// New event
		text = "新事件\n" + e.String()
		msg, err = b.SendMessage(text)
	} else if !oldMsg.event.Equal(e) {
		// update old event
		text = "事件更新\n" + oldMsg.event.Diff(e) + "\n" + e.String()
		msg, err = b.SendMessage(text)
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
		b.SendMessage(fmt.Sprintf("||--debug--\ngc: %d -> %d||", l, b.bucket.Len()))
	}
}
