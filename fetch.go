package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var target_url = "https://119dts.tncfd.gov.tw/DTS/caselist/html"

func init() {
	target_url = Getenv("TARGET_URL", target_url)
}

func fetch(filter func(Event) bool) (map[string]Event, error) {
	res, err := http.Get(target_url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("failed to fetch data: " + res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	results := map[string]Event{}

	doc.Find("tbody tr").Not(":first-child").Each(func(i int, s *goquery.Selection) {
		t, err := time.Parse("2006/01/02 15:04:05", s.Find(":nth-child(3)").Text())
		if err == nil {
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local)
		}

		e := Event{
			ID:       s.Find(":nth-child(2)").Text(),
			Time:     t,
			Type:     s.Find(":nth-child(4)").Text(),
			Location: s.Find(":nth-child(5)").Text(),
			Brigade:  strings.Split(s.Find(":nth-child(6)").Text(), ","),
			Status:   s.Find(":nth-child(7)").Text(),
		}

		if filter(e) {
			results[e.ID] = e
		}
	})

	return results, nil
}
