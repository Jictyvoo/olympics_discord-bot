package services

import (
	"sync"
	"testing"
)

type recorder struct{ got []int }

func (r *recorder) On(v int) { r.got = append(r.got, v) }

func TestSubject_Emit_FansOutToAllObservers(t *testing.T) {
	s := NewSubject[int]()
	a, b, c := &recorder{}, &recorder{}, &recorder{}
	s.Register(a)
	s.Register(b)
	s.Register(c)

	s.Emit(42)

	for _, obs := range []*recorder{a, b, c} {
		if len(obs.got) != 1 || obs.got[0] != 42 {
			t.Fatalf("observer got %v, want [42]", obs.got)
		}
	}
}

func TestSubject_Emit_NoObservers_NoPanic(t *testing.T) {
	NewSubject[int]().Emit(1)
}

type noopObserver struct{}

func (noopObserver) On(int) {}

// Run with -race to verify the snapshot-under-lock in Emit. Observers are
// stateless so the only shared state under test is the subject's own slice.
func TestSubject_ConcurrentRegisterAndEmit(t *testing.T) {
	s := NewSubject[int]()
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() { s.Register(noopObserver{}) })
		wg.Go(func() { s.Emit(1) })
	}
	wg.Wait()
}
