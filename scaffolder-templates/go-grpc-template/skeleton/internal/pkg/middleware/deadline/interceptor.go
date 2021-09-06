package deadline

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor returns a new unary server interceptor for error stack trace.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		switch ctx.Err() {
		case context.DeadlineExceeded:
			return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
		case context.Canceled:
			return nil, status.Error(codes.Canceled, "canceled")
		}
		resp, err := handler(ctx, req)
		switch ctx.Err() {
		case context.DeadlineExceeded:
			return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
		case context.Canceled:
			return nil, status.Error(codes.Canceled, "canceled")
		}
		return resp, err
	}
}

// TODO: implement StreamServerInterceptor
// // StreamServerInterceptor returns a new streaming server interceptor for error stack trace.
// func StreamServerInterceptor() grpc.StreamServerInterceptor {
// 	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
// 		return handler(srv, stream)
// 	}
// }
