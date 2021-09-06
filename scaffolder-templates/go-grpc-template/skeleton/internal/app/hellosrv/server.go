package hellosrv

import (
	"context"

	helloPb "${{values.go_mod_url}}/api/v1/hello"
	"${{values.go_mod_url}}/internal/pkg/service"
)

type HelloServer struct {
	helloPb.UnimplementedGreeterServer
	helloService *service.HelloService
}

func NewHelloServer() *HelloServer {
	hs := service.NewHelloService()
	return &HelloServer{
		helloService: hs,
	}
}

func (hs *HelloServer) SayHello(ctx context.Context, request *helloPb.HelloRequest) (*helloPb.HelloReply, error) {
	return hs.helloService.SayHello(ctx, request)
}
