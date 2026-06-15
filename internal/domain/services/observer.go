package services

import "sync"

// Observer reacts to an emitted value of type T.
type Observer[T any] interface {
	On(value T)
}

// Subject fans an emitted value out to every registered observer. Observers are
// held by value (the interface), so callers register them without taking an
// address; the subject keeps them alive for its lifetime.
type Subject[T any] struct {
	mu        sync.Mutex
	observers []Observer[T]
}

func NewSubject[T any]() *Subject[T] { return &Subject[T]{} }

func (s *Subject[T]) Register(observer Observer[T]) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observers = append(s.observers, observer)
}

// Emit delivers value to every registered observer. The observer slice is
// snapshotted under the lock so a slow observer doesn't block Register.
func (s *Subject[T]) Emit(value T) {
	s.mu.Lock()
	observers := make([]Observer[T], len(s.observers))
	copy(observers, s.observers)
	s.mu.Unlock()

	for _, observer := range observers {
		observer.On(value)
	}
}
