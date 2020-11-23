package defs

import "time"

//Order :
type Order struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	PrepTime time.Duration `json:"prepTime"`
}
