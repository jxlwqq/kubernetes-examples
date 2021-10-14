package main

import (
	calculatorv1 "github.com/jxlwqq/route-guide/api/protobuf/calculator"
	healthv1 "github.com/jxlwqq/route-guide/api/protobuf/health"
	"github.com/jxlwqq/route-guide/internal/calculator"
	"github.com/jxlwqq/route-guide/internal/health"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	PORT = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	calculatorServer := calculator.NewServer()
	healthServer := health.NewServer()

	calculatorv1.RegisterCalculatorServer(s, calculatorServer)
	healthv1.RegisterHealthServer(s, healthServer)

	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
