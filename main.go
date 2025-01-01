package main

import (
	"log"
	"time"

	"github.com/expr-lang/expr"
)

var (
	filterStr       = `!(Status == "已到達" || Status == "已到院" || Status == "返隊中" || Status == "已返隊") && (Type != "緊急救護" || len(Detachment) >= 2)`
	api_key         = ""
	api             = "https://api.telegram.org/bot"
	chat_id   int64 = 0
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
	api_key = Getenv("API_KEY", api_key)
	api = Getenv("API", api)
	chat_id = GetenvInt64("CHAT_ID", chat_id)

	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

func main() {
	bot := NewBot(api_key, chat_id)

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
