package queue

import "github.com/dispatch-simulator/internal/defs"

//AccessOrderQueue :
type AccessOrderQueue interface {
	Enqueue(a defs.Item)
	Dequeue() *defs.Item
}

//AccessDispatchQueue :
type AccessDispatchQueue interface {
	Enqueue(a defs.Item)
	Dequeue() *defs.Item
}
