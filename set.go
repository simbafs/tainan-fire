package main

type Set[T any] map[string]T

// A.Diff(B) returns A-B, B-A
func (A Set[T]) Diff(B Set[T]) ([]T, []T) {
	AB := []T{} // A-B
	BA := []T{} // B-A

	for k, v := range A {
		if _, ok := B[k]; !ok {
			AB = append(AB, v)
		}
	}

	for k, v := range B {
		if _, ok := A[k]; !ok {
			BA = append(BA, v)
		}
	}

	return AB, BA
}
