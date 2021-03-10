package hello_service

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"grpc/basic/prototypes"
)

type HelloService struct {
	prototypes.UnimplementedHelloServiceServer
}

func (x HelloService) SayHello(ctx context.Context, in *empty.Empty) (*prototypes.HelloMessage, error) {
	return &prototypes.HelloMessage{Value: "Hello"}, nil
}
