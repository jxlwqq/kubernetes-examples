package calculator

import (
	"context"
	calculatorv1 "github.com/jxlwqq/route-guide/api/protobuf/calculator"
)

type server struct {
	calculatorv1.UnimplementedCalculatorServer
}

func NewServer() calculatorv1.CalculatorServer {
	return &server{}
}

func (server) Add(ctx context.Context, req *calculatorv1.Request) (*calculatorv1.Response, error) {
	x := req.X
	y := req.Y
	res := x + y
	resp := calculatorv1.Response{
		Res: res,
		Err: "",
	}

	return &resp, nil
}
func (server) Subtract(ctx context.Context, req *calculatorv1.Request) (*calculatorv1.Response, error) {
	x := req.X
	y := req.Y
	res := x - y
	resp := calculatorv1.Response{
		Res: res,
		Err: "",
	}

	return &resp, nil
}
func (server) Multiply(ctx context.Context, req *calculatorv1.Request) (*calculatorv1.Response, error) {
	x := req.X
	y := req.Y
	res := x * y
	resp := calculatorv1.Response{
		Res: res,
		Err: "",
	}

	return &resp, nil
}
func (server) Divide(ctx context.Context, req *calculatorv1.Request) (*calculatorv1.Response, error) {
	x := req.X
	y := req.Y
	if y == 0 {
		e := "divisor can not be 0"
		resp := calculatorv1.Response{Err: e}

		return &resp, nil
	}
	res := x / y
	resp := calculatorv1.Response{
		Res: res,
		Err: "",
	}

	return &resp, nil
}
