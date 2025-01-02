package main

import (
	"fmt"
	"strings"
	"time"
)

type List []string

func (d List) String() string {
	return strings.Join(d, ",")
}

func (d List) Equal(other List) bool {
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

func (l List) Diff(other List) string {
	if l.Equal(other) {
		return ""
	}

	addition := []string{}
	deletion := []string{}
	// addition
	for _, v := range other {
		found := false
		for _, vv := range l {
			if v == vv {
				found = true
				break
			}
		}
		if !found {
			addition = append(addition, v)
		}
	}
	// deletion
	for _, v := range l {
		found := false
		for _, vv := range other {
			if v == vv {
				found = true
				break
			}
		}
		if !found {
			deletion = append(deletion, v)
		}
	}

	s := ""
	if len(addition) > 0 {
		s += fmt.Sprintf("新增: %s\n", strings.Join(addition, ","))
	}
	if len(deletion) > 0 {
		s += fmt.Sprintf("刪除: %s\n", strings.Join(deletion, ","))
	}

	return s
}

type Event struct {
	ID       string
	Time     time.Time
	Type     string
	Location string
	Brigade  List
	Status   string
}

const timeLayout = "2006/01/02 15:04:05"

func (e *Event) String() string {
	s := ""

	if len(e.Brigade) == 0 {
		s += fmt.Sprintf("`%s\n%s %s %s`", e.Time.Format(timeLayout), e.Type, e.Location, e.Status)
	} else {
		s += fmt.Sprintf("`%s\n%s %s %s\n%s`", e.Time.Format(timeLayout), e.Type, e.Location, e.Status, e.Brigade)
	}

	// debug //
	s += fmt.Sprintf("\n||---debug---\n%s||", e.ID)

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
		e.Brigade.Equal(New.Brigade) &&
		e.Status == New.Status
}

func (e *Event) Diff(other *Event) string {
	if e.Equal(other) {
		return ""
	}

	s := ""
	if e.Time != other.Time {
		s += fmt.Sprintf("時間: %s -> %s\n", e.Time.Format(timeLayout), other.Time.Format(timeLayout))
	}
	if e.Type != other.Type {
		s += fmt.Sprintf("類型: %s -> %s\n", e.Type, other.Type)
	}
	if e.Location != other.Location {
		s += fmt.Sprintf("地點: %s -> %s\n", e.Location, other.Location)
	}
	if !e.Brigade.Equal(other.Brigade) {
		s += fmt.Sprintf("分隊：\n%s", e.Brigade.Diff(other.Brigade))
	}
	if e.Status != other.Status {
		s += fmt.Sprintf("狀態: %s -> %s\n", e.Status, other.Status)
	}
	return s
}
