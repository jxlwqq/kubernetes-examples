package main

import (
	"context"
	pb "github.com/jxlwqq/route-guide/api/protobuf/calculator"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	PORT = ":50051"
)

type server struct {
	pb.UnimplementedCalculatorServer
}

func (server) Add(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	x := req.X
	y := req.Y
	res := x + y
	resp := pb.Response{
		Res: res,
		Err: "",
	}

	return &resp, nil
}
func (server) Subtract(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	x := req.X
	y := req.Y
	res := x - y
	resp := pb.Response{
		Res: res,
		Err: "",
	}

	return &resp, nil
}
func (server) Multiply(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	x := req.X
	y := req.Y
	res := x * y
	resp := pb.Response{
		Res: res,
		Err: "",
	}

	return &resp, nil
}
func (server) Divide(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	x := req.X
	y := req.Y
	if y == 0 {
		e := "divisor can not be 0"
		resp := pb.Response{Err: e}

		return &resp, nil
	}
	res := x / y
	resp := pb.Response{
		Res: res,
		Err: "",
	}

	return &resp, nil
}

func main() {
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCalculatorServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
