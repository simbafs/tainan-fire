package main

import (
	"log"
	"time"

	"github.com/expr-lang/expr"
)

var filterStr = `!(Status == "已到達" || Status == "已到院" || Status == "返隊中" || Status == "已返隊") && (Type != "緊急救護" || len(Detachment) >= 2)`

func init() {
	filterStr = Getenv("FILTER", filterStr)
}

func main() {
	h := Set[Event]{}

	filter := func(e Event) bool {
		// return !(e.Status == "已到達" || e.Status == "已到院" || e.Status == "返隊中" || e.Status == "已返隊")
		p, err := expr.Compile(filterStr, expr.Env(e))
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

	for {
		events, err := fetch(filter)
		if err != nil {
			log.Println(err)
		}

		_, newEvents := h.Diff(events)
		h = events

		s := ""
		for _, event := range newEvents {
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
