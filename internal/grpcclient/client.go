package grpcclient

import (
	"google.golang.org/grpc"
	"homeWorkAdvancedPart2_REST_service/internal/pb"
)

type Client struct {
	pb.DeliveryServiceClient
}

func NewClient(conn grpc.ClientConnInterface) *Client {
	c := pb.NewDeliveryServiceClient(conn)
	return &Client{c}
}
