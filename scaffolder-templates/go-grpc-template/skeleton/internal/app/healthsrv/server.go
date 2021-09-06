package healthsrv

import (
	"context"

	"google.golang.org/grpc/codes"
	hc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type server struct {
}

// Check returns if the service is available.
func (s *server) Check(c context.Context, r *hc.HealthCheckRequest) (*hc.HealthCheckResponse, error) {
	// Add service specific health check logic here.
	// i.e. Make sure if db is ready.
	return &hc.HealthCheckResponse{
		Status: hc.HealthCheckResponse_SERVING,
	}, nil
}

// Watch is used to check health using streaming. Not using it for now.
func (s *server) Watch(*hc.HealthCheckRequest, hc.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "not implemented")
}

// New initializes health server
func New() hc.HealthServer {
	return &server{}
}
