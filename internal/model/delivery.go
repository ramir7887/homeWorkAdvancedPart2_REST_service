package model

import "time"

type Delivery struct {
	ID           int
	OrderID      int
	DeliveryDate time.Time
	Complete     bool
	Address      string
}
