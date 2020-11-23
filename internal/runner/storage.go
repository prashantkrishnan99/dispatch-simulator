package runner

import (
	"sync"
)

//Storage : Data structure defining storage
type Storage struct {
	mu     sync.RWMutex
	orders map[string]interface{}
}

//NewStorage :
func NewStorage() *Storage {
	return &Storage{orders: make(map[string]interface{})}
}

//Insert : Insert data into the storage
func (s *Storage) Insert(key string, data interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[key] = data
	return
}

//Delete : Deletes an entry from storage
func (s *Storage) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.orders, key)
}

//Get :
func (s *Storage) Get(key string) interface{} {
	s.mu.RLock()
	if order, ok := s.orders[key]; ok {
		s.mu.RUnlock()
		return order
	}
	s.mu.RUnlock()
	return nil
}

//Dump : Dumps the storage
func (s *Storage) Dump() interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.orders
}

//Flush : Flushes the storage
func (s *Storage) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders = make(map[string]interface{})
}

//IsEmpty : Checks if the storage is Empty or not
func (s *Storage) IsEmpty() bool {
	if len(s.orders) == 0 {
		return true
	}
	return false
}
