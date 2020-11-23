package defs

import "github.com/cheekybits/genny/generic"

const (
	//Matched :
	Matched = iota
	//Fifo :
	Fifo
)

//ORDERREADY :
const ORDERREADY = "orderready_"

//DISPATCHREADY :
const DISPATCHREADY = "dispatchready_"

//Item : queue interface type
type Item generic.Type
