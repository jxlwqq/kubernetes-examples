package main

import (
	"context"
	"github.com/gin-gonic/gin"
	calculatorv1 "github.com/jxlwqq/route-guide/api/protobuf/calculator"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	ADDRESS = "localhost"
	PORT    = ":50051"
)

var address string

func main() {

	viper.AutomaticEnv()
	viper.SetDefault("CALCULATOR_SVC", ADDRESS)
	address = viper.GetString("CALCULATOR_SVC")
	r := gin.Default()
	r.GET("/add/:x/:y", Add)
	r.GET("/subtract/:x/:y", Subtract)
	r.GET("/multiply/:x/:y", Multiply)
	r.GET("/divide/:x/:y", Divide)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("server serve failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("server is shutting down ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server is exiting")
}

func Multiply(c *gin.Context) {
	x, _ := strconv.ParseFloat(c.Param("x"), 32)
	y, _ := strconv.ParseFloat(c.Param("y"), 32)

	res, err := calculate(float32(x), "*", float32(y))
	c.JSON(http.StatusOK, gin.H{
		"res": res,
		"err": err,
	})
}

func Divide(c *gin.Context) {
	x, _ := strconv.ParseFloat(c.Param("x"), 32)
	y, _ := strconv.ParseFloat(c.Param("y"), 32)

	res, err := calculate(float32(x), "/", float32(y))
	c.JSON(http.StatusOK, gin.H{
		"res": res,
		"err": err,
	})
}

func Subtract(c *gin.Context) {
	x, _ := strconv.ParseFloat(c.Param("x"), 32)
	y, _ := strconv.ParseFloat(c.Param("y"), 32)

	res, err := calculate(float32(x), "-", float32(y))
	c.JSON(http.StatusOK, gin.H{
		"res": res,
		"err": err,
	})
}

func Add(c *gin.Context) {
	x, _ := strconv.ParseFloat(c.Param("x"), 32)
	y, _ := strconv.ParseFloat(c.Param("y"), 32)

	res, err := calculate(float32(x), "+", float32(y))
	c.JSON(http.StatusOK, gin.H{
		"res": res,
		"err": err,
	})
}

func calculate(x float32, operator string, y float32) (float32, string) {
	conn, err := grpc.Dial(address+PORT, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := calculatorv1.NewCalculatorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &calculatorv1.Request{X: x, Y: y}

	switch operator {
	case "+":
		resp, _ := c.Add(ctx, req)
		return resp.Res, resp.Err
	case "-":
		resp, _ := c.Subtract(ctx, req)
		return resp.Res, resp.Err
	case "*":
		resp, _ := c.Multiply(ctx, req)
		return resp.Res, resp.Err
	case "/":
		resp, _ := c.Divide(ctx, req)
		return resp.Res, resp.Err
	}

	return 0, ""
}
