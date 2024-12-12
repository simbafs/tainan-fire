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
	events := NewEvents(bot.SendEvent)

	for {
		e, err := fetch(filter)
		if err != nil {
			log.Println(err)
		}

		err = events.Update(e)

		// sort by time
		// for i := 0; i < len(sortedEvents); i++ {
		// 	for j := i + 1; j < len(sortedEvents); j++ {
		// 		if sortedEvents[i].Time.After(sortedEvents[j].Time) {
		// 			sortedEvents[i], sortedEvents[j] = sortedEvents[j], sortedEvents[i]
		// 		}
		// 	}
		// }

		time.Sleep(30 * time.Second)
	}
}
