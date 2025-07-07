package main

import (
	"log"
	"time"

	"github.com/expr-lang/expr"
)

var (
	filterStr       = `!(Status == "已到達" || Status == "已到院" || Status == "返隊中" || Status == "已返隊") && (Type != "緊急救護" || len(Brigade) >= 2)`
	apiKEY          = ""
	api             = "https://api.telegram.org/bot"
	chatID    int64 = 0
)

func filter(e Event) bool {
	p, err := expr.Compile(filterStr, expr.Env(e), expr.AsBool())
	if err != nil {
		log.Println(err)
		return false
	}

	r, err := expr.Run(p, e)
	if err != nil {
		log.Println(err)
		return false
	}

	return r.(bool)
}

func init() {
	filterStr = Getenv("FILTER", filterStr)
	apiKEY = Getenv("API_KEY", apiKEY)
	api = Getenv("API", api)
	chatID = GetenvInt64("CHAT_ID", chatID)

	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

func main() {
	bot := NewBot(WithAPIKey(apiKEY), WithChatID(chatID))

	for {
		events, err := fetch(filter)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, event := range events {
			if err := bot.SendEvent(&event); err != nil {
				log.Println(err)
			}
		}

		bot.GC()

		time.Sleep(10 * time.Second)
	}
}
