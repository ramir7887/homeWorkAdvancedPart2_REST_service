package server

import (
	"context"
	"encoding/json"
	"homeWorkAdvancedPart2_REST_service/internal/genserve"
	"homeWorkAdvancedPart2_REST_service/internal/grpcclient"
	"homeWorkAdvancedPart2_REST_service/internal/model"
	"homeWorkAdvancedPart2_REST_service/internal/pb"
	"homeWorkAdvancedPart2_REST_service/internal/store"
	"log"
	"net/http"
)

type server struct {
	store store.Store
	gc    *grpcclient.Client
	genserve.ServerInterface
}

func NewServer(store store.Store, gc *grpcclient.Client) *server {
	return &server{
		store: store,
		gc:    gc,
	}
}
func (s *server) GetOrders(w http.ResponseWriter, r *http.Request, params genserve.GetOrdersParams) { // Done
	limit := 0
	offset := 0
	if params.Limit != nil {
		limit = *params.Limit
	}
	if params.Offset != nil {
		offset = *params.Offset
	}
	orders, err := s.store.Order().GetOrders(limit, offset)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Извините, можно было бы написать мапперы
	outOrders := make([]genserve.Order, 0, len(orders))
	for _, v := range orders {
		order := v
		deliveryDate := order.Delivery.DeliveryDate.String()
		outOrder := genserve.Order{
			Id:          &order.ID,
			Address:     order.Address,
			Description: &order.Description,
			Delivery: &genserve.Delivery{
				Id:           &order.Delivery.ID,
				Address:      &order.Delivery.Address,
				OrderId:      &order.Delivery.OrderID,
				Complete:     &order.Delivery.Complete,
				DeliveryDate: &deliveryDate,
			},
			Products: func(o model.Order) []genserve.Product {
				products := make([]genserve.Product, 0, len(o.Products))
				for _, op := range o.Products {
					opId := op.ID
					gp := genserve.Product{
						Id:    &opId,
						Name:  op.Name,
						Price: float32(op.Price),
					}
					products = append(products, gp)
				}
				return products
			}(order),
		}
		outOrders = append(outOrders, outOrder)
	}

	jsonData, err := json.Marshal(outOrders)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (s *server) AddOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("AddOrder")
	ctx := context.Background()

	inOrder := &genserve.CreateOrder{}
	if err := json.NewDecoder(r.Body).Decode(inOrder); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	order := model.Order{
		Address:     inOrder.Address,
		Description: *inOrder.Description,
		Products:    make([]model.Product, 0, len(inOrder.Products)),
	}
	for _, v := range inOrder.Products {
		product, err := s.store.Product().GetProduct(v)
		if err != nil {
			log.Println("create product for order: ", err)
			continue
		}
		order.Products = append(order.Products, *product)
	}

	deliveryOrder := pb.Order{Id: int32(s.store.Order().NextOrderID()), Address: inOrder.Address}
	deliveryId, err := s.gc.CreateDelivery(ctx, &deliveryOrder)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	delivery, err := s.gc.GetDelivery(ctx, deliveryId)
	if err != nil || delivery == nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	order.Delivery = model.Delivery{
		ID:           int(delivery.Id),
		Address:      delivery.Address,
		OrderID:      int(delivery.OrderId),
		Complete:     delivery.Complete,
		DeliveryDate: delivery.DeliveryDate.AsTime(),
	}

	id := s.store.Order().CreateOrder(&order)
	m := map[string]int{"id": id}
	jsonData, err := json.Marshal(m)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
} //Done

func (s *server) DeleteOrder(w http.ResponseWriter, r *http.Request, orderId int) {
	log.Println("DeleteOrder")
	ctx := context.Background()
	order, err := s.store.Order().GetOrder(orderId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	complete, err := s.gc.DeleteDelivery(ctx, &pb.DeliveryID{Value: int32(order.ID)})
	if err != nil || !complete.Value {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := s.store.Order().DeleteOrder(orderId); err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *server) GetOrderById(w http.ResponseWriter, r *http.Request, orderId int) {
	log.Println("GetOrderById", orderId)
	order, err := s.store.Order().GetOrder(int(orderId))
	if err != nil || order == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	jsonData, err := json.Marshal(order)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (s *server) GetProducts(w http.ResponseWriter, r *http.Request, params genserve.GetProductsParams) {
	log.Println("GetProducts", params)
	limit := 0
	offset := 0
	if params.Limit != nil {
		limit = *params.Limit
	}
	if params.Offset != nil {
		offset = *params.Offset
	}

	products, err := s.store.Product().GetProducts(limit, offset)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	outProducts := make([]genserve.Product, 0, len(products))
	for _, v := range products {
		product := v
		outProduct := genserve.Product{
			Id:          &product.ID,
			Name:        product.Name,
			Price:       float32(product.Price),
			Description: &product.Description,
		}
		outProducts = append(outProducts, outProduct)
	}

	jsonData, err := json.Marshal(outProducts)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
} // Done

func (s *server) AddProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("AddProduct")
	inProduct := &genserve.Product{}
	if err := json.NewDecoder(r.Body).Decode(inProduct); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	product := model.Product{
		ID:          *inProduct.Id,
		Name:        inProduct.Name,
		Price:       float64(inProduct.Price),
		Description: *inProduct.Description,
	}

	id := s.store.Product().CreateProduct(&product)
	m := map[string]int{"id": id}
	jsonData, err := json.Marshal(m)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

} // Done

func (s *server) DeleteProduct(w http.ResponseWriter, r *http.Request, productId int) {
	log.Println("DeleteProduct", productId)
	if err := s.store.Product().DeleteProduct(productId); err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
} // Done

func (s *server) GetProductById(w http.ResponseWriter, r *http.Request, productId int) {
	log.Println("GetProductById", productId)
	product, err := s.store.Product().GetProduct(productId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	outProduct := genserve.Product{
		Id:          &product.ID,
		Name:        product.Name,
		Price:       float32(product.Price),
		Description: &product.Description,
	}
	jsonData, err := json.Marshal(outProduct)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
} // Done
