package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	helloPb "${{values.go_mod_url}}/api/v1/hello"
	"${{values.go_mod_url}}/internal/app/hellosrv"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"${{values.go_mod_url}}/internal/app/healthsrv"
	"github.com/sirupsen/logrus"

	middleware_deadline "${{values.go_mod_url}}/internal/pkg/middleware/deadline"
	middleware_stacktrace "${{values.go_mod_url}}/internal/pkg/middleware/stacktrace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	hc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", "6565"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		os.Interrupt,
		os.Kill,
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
		syscall.SIGTERM,
	)

	s := newGrpcServer()
	logrus.SetLevel(logrus.InfoLevel)

	helloPb.RegisterGreeterServer(s, hellosrv.NewHelloServer())
	hc.RegisterHealthServer(s, healthsrv.New())

	go func() {
		<-signalChan
		s.GracefulStop()
	}()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
		signalChan <- syscall.SIGINT
	}
	s.GracefulStop()
}

func newGrpcServer() *grpc.Server {
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
	logrusOpts := []grpc_logrus.Option{
		grpc_logrus.WithDecider(func(methodFullName string, err error) bool {
			if methodFullName == "/grpc.health.v1.Health/Check" {
				return false
			}
			return true
		}),
	}
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)

	payloadLoggingDecider := func(ctx context.Context, methodFullName string, servingObject interface{}) bool {
		if methodFullName == "/grpc.health.v1.Health/Check" {
			return false
		}
		return true
	}

	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			return status.Errorf(codes.Unknown, "panic triggered: %+v: %+v", p, string(debug.Stack()))
		}),
	}

	return grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logrusEntry, logrusOpts...),
			middleware_stacktrace.UnaryServerInterceptor(),
			grpc_logrus.PayloadUnaryServerInterceptor(logrusEntry, payloadLoggingDecider),
			middleware_deadline.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.StreamServerInterceptor(logrusEntry, logrusOpts...),
			middleware_stacktrace.StreamServerInterceptor(),
			grpc_logrus.PayloadStreamServerInterceptor(logrusEntry, payloadLoggingDecider),
			grpc_recovery.StreamServerInterceptor(recoveryOpts...),
		),
	)
}
