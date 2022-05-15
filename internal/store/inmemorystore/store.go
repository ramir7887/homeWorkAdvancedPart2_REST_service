package inmemorystore

import (
	"homeWorkAdvancedPart2_REST_service/internal/model"
	"homeWorkAdvancedPart2_REST_service/internal/store"
	"time"
)

type Store struct {
	orders            []model.Order
	products          []model.Product
	orderRepository   store.OrderRepository
	productRepository store.ProductRepository
}

func NewStore() *Store {
	products := []model.Product{
		{ID: 1, Name: "Shampoo", Price: 430, Description: "Nice cool"},
		{ID: 2, Name: "Banana", Price: 96.45, Description: "Nice cool banana"},
		{ID: 3, Name: "Cola", Price: 46.45, Description: "Nice cool Cola"},
		{ID: 4, Name: "Pepsi", Price: 46.45, Description: "Nice cool Pepsi"},
		{ID: 5, Name: "Tomato", Price: 16.45, Description: "Nice cool Tomato"},
		{ID: 6, Name: "Cake", Price: 86.45, Description: "Nice cool Cake"},
		{ID: 7, Name: "WaterMelon", Price: 76.45, Description: "Nice cool WaterMelon"},
		{ID: 8, Name: "Burger", Price: 246.45, Description: "Nice cool Burger"},
		{ID: 9, Name: "Vodka", Price: 146.45, Description: "Nice cool Vodka"},
	}

	orders := []model.Order{
		{ID: 1, Description: "To Tundra", Address: "Norilsk, Lenin street 16, 15", Products: products[:3], Delivery: model.Delivery{ID: 1, DeliveryDate: time.Now().Add(24 * 3 * time.Hour), Address: "Norilsk, Lenin street 16, 15", Complete: false, OrderID: 1}},
		{ID: 2, Description: "Очень аккуратно. Хрупкий товар", Address: "", Products: products[:4], Delivery: model.Delivery{ID: 1, DeliveryDate: time.Now().Add(24 * 3 * time.Hour), Address: "Norilsk, Lenin street 16, 15", Complete: false, OrderID: 1}},
		{ID: 3, Description: "Для друга", Address: "", Products: products[2:5], Delivery: model.Delivery{ID: 1, DeliveryDate: time.Now().Add(24 * 3 * time.Hour), Address: "Norilsk, Lenin street 16, 15", Complete: false, OrderID: 1}},
		{ID: 4, Description: "Срочная доставка", Address: "", Products: products[3:5], Delivery: model.Delivery{ID: 1, DeliveryDate: time.Now().Add(24 * 3 * time.Hour), Address: "Norilsk, Lenin street 16, 15", Complete: false, OrderID: 1}},
		{ID: 5, Description: "Звонить в домофон", Address: "", Products: products[4:5], Delivery: model.Delivery{ID: 1, DeliveryDate: time.Now().Add(24 * 3 * time.Hour), Address: "Norilsk, Lenin street 16, 15", Complete: false, OrderID: 1}},
		{ID: 6, Description: "Не кантовать", Address: "", Products: products[1:3], Delivery: model.Delivery{ID: 1, DeliveryDate: time.Now().Add(24 * 3 * time.Hour), Address: "Norilsk, Lenin street 16, 15", Complete: false, OrderID: 1}},
		{ID: 8, Description: "Обнимать", Address: "", Products: products[:], Delivery: model.Delivery{ID: 1, DeliveryDate: time.Now().Add(24 * 3 * time.Hour), Address: "Norilsk, Lenin street 16, 15", Complete: false, OrderID: 1}},
	}

	return &Store{
		orders:   orders,
		products: products,
	}
}

func (s *Store) Order() store.OrderRepository {
	if s.orderRepository != nil {
		return s.orderRepository
	}

	s.orderRepository = &OrderRepository{
		store: s,
	}
	return s.orderRepository
}

func (s *Store) Product() store.ProductRepository {
	if s.productRepository != nil {
		return s.productRepository
	}

	s.productRepository = &ProductRepository{
		store: s,
	}
	return s.productRepository
}
