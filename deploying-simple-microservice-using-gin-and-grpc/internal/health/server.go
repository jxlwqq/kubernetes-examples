package health

import (
	"context"
	healthv1 "github.com/jxlwqq/route-guide/api/protobuf/health"
	"log"
)

type server struct {
	healthv1.UnimplementedHealthServer
}

func (s server) Check(ctx context.Context, request *healthv1.HealthCheckRequest) (*healthv1.HealthCheckResponse, error) {
	log.Printf("Received Check")
	return &healthv1.HealthCheckResponse{
		Status: healthv1.HealthCheckResponse_SERVING,
	}, nil
}

func (s server) Watch(request *healthv1.HealthCheckRequest, watchServer healthv1.Health_WatchServer) error {
	log.Printf("Received Watch")
	return watchServer.Send(&healthv1.HealthCheckResponse{
		Status: healthv1.HealthCheckResponse_SERVING,
	})
}

func NewServer() healthv1.HealthServer {
	return &server{}
}
