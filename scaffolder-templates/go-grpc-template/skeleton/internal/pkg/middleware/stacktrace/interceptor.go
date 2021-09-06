package stacktrace

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor returns a new unary server interceptor for error stack trace.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			stacktrace := fmt.Sprintf("%+v", err)
			field := logrus.Fields{
				"stacktrace": stacktrace,
			}
			ctxlogrus.AddFields(ctx, field)
		}
		return resp, err
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for error stack trace.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, stream)
		if err != nil {
			stacktrace := fmt.Sprintf("%+v", err)
			field := logrus.Fields{
				"stacktrace": stacktrace,
			}
			// TODO: check whether it works
			ctxlogrus.AddFields(stream.Context(), field)
		}
		return err
	}
}
