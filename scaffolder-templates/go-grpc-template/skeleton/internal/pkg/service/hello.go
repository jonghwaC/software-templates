package service

import (
	"context"

	helloPb "${{values.go_mod_url}}/api/v1/hello"
)

type HelloService struct {
}

func NewHelloService() *HelloService {
	return &HelloService{}
}

func (hs *HelloService) SayHello(ctx context.Context, request *helloPb.HelloRequest) (*helloPb.HelloReply, error) {
	response := &helloPb.HelloReply{
		Message: "Hello " + request.Name,
	}
	return response, nil
}
