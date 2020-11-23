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
