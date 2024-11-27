package main

type History[T any] map[string]T

func (h History[T]) Check(updates map[string]T) []T {
	var newEvents []T

	for k, v := range updates {
		if _, ok := h[k]; !ok {
			newEvents = append(newEvents, v)
		}
		h[k] = v
	}

	return newEvents
}
