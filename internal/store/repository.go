package store

import "homeWorkAdvancedPart2_REST_service/internal/model"

type ProductRepository interface {
	CreateProduct(*model.Product) int
	GetProduct(int) (*model.Product, error)
	GetProducts(int, int) ([]model.Product, error)
	DeleteProduct(int) error
}

type OrderRepository interface {
	CreateOrder(*model.Order) int
	GetOrder(int) (*model.Order, error)
	GetOrders(int, int) ([]model.Order, error)
	DeleteOrder(int) error
	NextOrderID() int
}

type ProductOrderRepository interface {
	ProductRepository
	OrderRepository
}
