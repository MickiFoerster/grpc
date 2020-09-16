package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/MickiFoerster/grpc/greet/go"
	"google.golang.org/grpc"
)

type serviceServer struct {
}

func (*serviceServer) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet() was invoked with %v\n", req)
	firstname := req.GetGreeting().GetFirstName()
	result := "Hello " + firstname
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Go server starts ...")

	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("listen failed: %v\n", err)
	}

	server := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(server, &serviceServer{})

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Serve() failed: %v\n", err)
	}
}
