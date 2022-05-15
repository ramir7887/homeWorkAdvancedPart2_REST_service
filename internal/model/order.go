package model

type Order struct {
	ID          int
	Description string
	Address     string
	Products    []Product
	Delivery    Delivery
}
