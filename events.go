package main

import (
	"fmt"
	"strings"
	"time"
)

type Detachment []string

func (d Detachment) String() string {
	return strings.Join(d, ",")
}

func (d Detachment) Equal(other Detachment) bool {
	if len(d) != len(other) {
		return false
	}
	for i, v := range d {
		if v != other[i] {
			return false
		}
	}
	return true
}

type Event struct {
	ID         string
	Time       time.Time
	Type       string
	Location   string
	Detachment Detachment
	Status     string
}

const timeLayout = "2006/01/02 15:04:05"

func (e *Event) String() string {
	s := ""

	if len(e.Detachment) == 0 {
		s += fmt.Sprintf("%s\n%s %s %s", e.Time.Format(timeLayout), e.Type, e.Location, e.Status)
	} else {
		s += fmt.Sprintf("%s\n%s %s %s\n%s", e.Time.Format(timeLayout), e.Type, e.Location, e.Status, e.Detachment)
	}

	// debug //
	s += fmt.Sprintf("\n---debug---\n%s", e.ID)

	return s
}

func (e *Event) Equal(New *Event) bool {
	if e == nil || New == nil {
		return false
	}
	return e.ID == New.ID &&
		e.Time == New.Time &&
		e.Type == New.Type &&
		e.Location == New.Location &&
		e.Detachment.Equal(New.Detachment) &&
		e.Status == New.Status
}

func (e *Event) Diff(New *Event) string {
	if e.Equal(New) {
		return ""
	}

	s := ""
	if e.Time != New.Time {
		s += fmt.Sprintf("Time: %s -> %s\n", e.Time.Format(timeLayout), New.Time.Format(timeLayout))
	}
	if e.Type != New.Type {
		s += fmt.Sprintf("Type: %s -> %s\n", e.Type, New.Type)
	}
	if e.Location != New.Location {
		s += fmt.Sprintf("Location: %s -> %s\n", e.Location, New.Location)
	}
	if !e.Detachment.Equal(New.Detachment) {
		s += fmt.Sprintf("Detachment: %s -> %s\n", e.Detachment, New.Detachment)
	}
	if e.Status != New.Status {
		s += fmt.Sprintf("Status: %s -> %s\n", e.Status, New.Status)
	}
	return s
}
