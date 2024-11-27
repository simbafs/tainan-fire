package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const timeLayout = "2006/01/02 15:04:05"

type Detachment []string

func (d Detachment) String() string {
	return strings.Join(d, ",")
}

type Event struct {
	ID         string
	Time       time.Time
	Type       string
	Location   string
	Detachment Detachment
	Status     string
}

func (e Event) String() string {
	if e.Detachment == nil {
		return fmt.Sprintf("%s %s %s %s", e.Time.Format(timeLayout), e.Type, e.Location, e.Status)
	}
	return fmt.Sprintf("%s %s %s %s %s", e.Time.Format(timeLayout), e.Type, e.Location, e.Detachment, e.Status)
}

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
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
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
		id := s.Find(":nth-child(2)").Text()

		e := Event{
			ID:         id,
			Time:       t,
			Type:       s.Find(":nth-child(4)").Text(),
			Location:   s.Find(":nth-child(5)").Text(),
			Detachment: strings.Split(s.Find(":nth-child(6)").Text(), ","),
			Status:     s.Find(":nth-child(7)").Text(),
		}

		if filter(e) {
			results[id] = e
		}
	})
	return results, nil
}
