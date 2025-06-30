package interpreter

import (
	"iter"
	"maps"
)

// This is an intentionally bare-bones implementation of a hash set
// that we only need once to keep track of bound variables
type set[T comparable] map[T]struct{}

func (s set[T]) Add(value T) {
	s[value] = struct{}{}
}

func (s set[T]) Has(value T) bool {
	_, ok := s[value]
	return ok
}

func (s set[T]) Values() iter.Seq[T] {
	return maps.Keys(s)
}
