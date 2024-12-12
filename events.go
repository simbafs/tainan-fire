package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

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

	Digest string

	PrevMessage *gotgbot.Message
}

const timeLayout = "2006/01/02 15:04:05"

func (e Event) String() string {
	if len(e.Detachment) == 0 {
		return fmt.Sprintf("%s\n%s %s %s\n---debug---\n%s %s", e.Time.Format(timeLayout), e.Type, e.Location, e.Status, e.ID, e.Digest)
	}
	return fmt.Sprintf("%s\n%s %s %s\n%s\n---debug---\n%s %s", e.Time.Format(timeLayout), e.Type, e.Location, e.Status, e.Detachment, e.ID, e.Digest)
}

type OnUpdate = func(*gotgbot.Message, *Event) error

type Events struct {
	events   map[string]Event // ID -> Event
	onUpdate OnUpdate
}

func NewEvents(onUpdate OnUpdate) *Events {
	return &Events{
		events:   map[string]Event{},
		onUpdate: onUpdate,
	}
}

// Update updates the events and returns new events and updated events
func (e *Events) Update(events []Event) error {
	err := NewErrors()

	idMap := map[string]struct{}{}

	for _, event := range events {
		idMap[event.ID] = struct{}{}
		if p, ok := e.events[event.ID]; !ok || p.Digest != event.Digest {
			err.Append(e.onUpdate(p.PrevMessage, &event))
		}

		e.events[event.ID] = event
	}

	for k := range e.events {
		if _, ok := idMap[k]; !ok {
			delete(e.events, k)
		}
	}

	return err
}
