package defs

//Store : Storage Interface
type Store interface {
	Insert(key string, data interface{})
	Delete(key string)
	Get(key string) interface{}
	Flush()
	Dump() interface{}
	IsEmpty() bool
}

//QueueStore : Storage Interface
type QueueStore interface {
	Enqueue(t Item)
	Dequeue() *Item
	Front() *Item
	Size() int
	IsEmpty() bool
}

//Stats :
type Stats interface {
	IncrOrdersProcessed()
	IncrTotalTime(int)
	CalculateAverage()
	GetTotalOrdersProcessed() int
	GetTotalTime() int
	GetAVerageTime() int
}
