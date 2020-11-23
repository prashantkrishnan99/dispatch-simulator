package defs

import "time"

//Order :
type Order struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	PrepTime time.Duration `json:"prepTime"`
}

//Courier :
type Courier struct {
	DispatchID string `json:"dispatch_id"`
}

//Dispatch :
type Dispatch struct {
	OrderID    string `json:"order_id"`
	DispatchID string `json:"dispatch_id"`
	Algo       int    `json:"algo"`
}
