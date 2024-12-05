package main

import (
	"strings"
	"time"
)

type Event struct {
	ID         string
	Time       time.Time
	Type       string
	Location   string
	Detachment Detachment
	Status     string

	Digest string

	MsgID int
}
type Detachment []string

func (d Detachment) String() string {
	return strings.Join(d, ",")
}

type Events struct {
	events map[string]Event // ID -> Event
}

func NewEvents() *Events {
	return &Events{
		events: map[string]Event{},
	}
}

// Update updates the events and returns new events and updated events
func (e *Events) Update(events []Event) ([]Event, []Event) {
	newEvents := []Event{}
	updateEvents := []Event{}

	for _, event := range events {
		if p, ok := e.events[event.ID]; ok {
			if p.Digest != event.Digest {
				updateEvents = append(updateEvents, event)
			}
		} else {
			newEvents = append(newEvents, event)
		}

		e.events[event.ID] = event
	}
	// TODO: remove old events

	return newEvents, updateEvents
}
