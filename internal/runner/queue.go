package runner

import (
	"sync"

	"github.com/dispatch-simulator/internal/defs"
)

//Queue :
type Queue struct {
	items []defs.Item
	lock  sync.RWMutex
}

//NewQueue :
func NewQueue() *Queue {
	return &Queue{
		items: []defs.Item{},
	}
}

//Enqueue :
func (s *Queue) Enqueue(t defs.Item) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.items = append(s.items, t)
}

//Dequeue :
func (s *Queue) Dequeue() *defs.Item {
	s.lock.Lock()
	defer s.lock.Unlock()
	if !s.IsEmpty() {
		item := s.items[0]
		s.items = s.items[1:len(s.items)]
		return &item
	}
	return nil
}

//Front :
func (s *Queue) Front() *defs.Item {
	s.lock.RLock()
	defer s.lock.RUnlock()
	item := s.items[0]
	return &item
}

//IsEmpty :
func (s *Queue) IsEmpty() bool {
	return len(s.items) == 0
}

//Size :
func (s *Queue) Size() int {
	return len(s.items)
}
