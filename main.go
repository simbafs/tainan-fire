package main

import (
	"log"
	"time"
)

func main() {
	h := History[Event]{}

	filter := func(e Event) bool {
		return !(e.Status == "已到達" || e.Status == "已到院" || e.Status == "返隊中" || e.Status == "已返隊")
	}

	for {
		events, err := fetch(filter)
		if err != nil {
			log.Println(err)
		}

		s := ""
		for _, event := range h.Check(events) {
			log.Println(event)
			s += event.String() + "\n"
		}

		if s != "" {
			s = "事件更新：\n" + s
			if err := SendMessage(s); err != nil {
				log.Println(err)
			}
		}

		time.Sleep(5 * time.Minute)
	}
}
