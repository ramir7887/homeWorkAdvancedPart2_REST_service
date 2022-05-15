package inmemorystore

import (
	"homeWorkAdvancedPart2_REST_service/internal/model"
	"log"
	"sync"
)

var (
	rnfError = RecordNotFoundError{}
)

type RecordNotFoundError struct {
}

func (r RecordNotFoundError) Error() string {
	return "record not found"
}

type OrderRepository struct {
	store *Store
	m     sync.Mutex
}

func (r *OrderRepository) CreateOrder(order *model.Order) int {
	id := r.NextOrderID()
	order.ID = id
	r.m.Lock()
	defer r.m.Unlock()
	r.store.orders = append(r.store.orders, *order)

	return id
}

func (r *OrderRepository) GetOrder(id int) (*model.Order, error) {
	r.m.Lock()
	defer r.m.Unlock()
	for _, val := range r.store.orders {
		if val.ID == id {
			order := val
			return &order, nil
		}
	}
	return nil, rnfError
}

func (r *OrderRepository) GetOrders(limit int, offset int) ([]model.Order, error) {
	if offset == 0 && limit == 0 {
		return r.store.orders, nil
	}
	if offset+limit >= len(r.store.orders) {
		return r.store.orders, nil
	}
	r.m.Lock()
	defer r.m.Unlock()
	return r.store.orders[offset : offset+limit], nil
}

func (r *OrderRepository) DeleteOrder(id int) error {
	index := -1
	for idx, val := range r.store.orders {
		if val.ID == id {
			index = idx
		}
	}
	if index == -1 {
		return rnfError
	}
	r.m.Lock()
	defer r.m.Unlock()
	r.store.orders = append(r.store.orders[:index], r.store.orders[index+1:]...)
	log.Printf("Record count: %v", len(r.store.orders))
	return nil
}

func (r *OrderRepository) lastOrderID() int {
	id := 0
	r.m.Lock()
	defer r.m.Unlock()
	for _, val := range r.store.orders {
		if val.ID > id {
			id = val.ID
		}
	}
	return id
}

func (r *OrderRepository) NextOrderID() int {
	return r.lastOrderID() + 1
}

type ProductRepository struct {
	store *Store
}

func (r *ProductRepository) CreateProduct(product *model.Product) int {
	id := r.nextProductID()
	product.ID = id
	r.store.products = append(r.store.products, *product)

	return id
}

func (r *ProductRepository) GetProduct(id int) (*model.Product, error) {
	for _, val := range r.store.products {
		if val.ID == id {
			product := val
			return &product, nil
		}
	}
	return nil, rnfError
}

func (r *ProductRepository) GetProducts(limit int, offset int) ([]model.Product, error) {
	if offset == 0 && limit == 0 {
		return r.store.products, nil
	}
	if offset+limit >= len(r.store.products) {
		return r.store.products, nil
	}
	return r.store.products[offset : offset+limit], nil

}

func (r *ProductRepository) DeleteProduct(id int) error {
	index := -1
	for idx, val := range r.store.products {
		if val.ID == id {
			index = idx
		}
	}
	if index == -1 {
		return rnfError
	}
	r.store.products = append(r.store.products[:index], r.store.products[index+1:]...)
	log.Printf("Record product count: %v", len(r.store.products))
	return nil
}

func (r *ProductRepository) lastProductID() int {
	id := 0
	for _, val := range r.store.products {
		if val.ID > id {
			id = val.ID
		}
	}
	return id
}

func (r *ProductRepository) nextProductID() int {
	return r.lastProductID() + 1
}
