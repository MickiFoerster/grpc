package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"

	greetwithdeadlinepb "github.com/MickiFoerster/grpc/greet-with-deadline"
)

type serviceServer struct {
}

func (*serviceServer) GreetWithDeadline(ctx context.Context, req *greetwithdeadlinepb.GreetWithDeadlineRequest) (*greetwithdeadlinepb.GreetWithDeadlineResponse, error) {
	fmt.Printf("Greet() was invoked with %v\n", req)

	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			log.Printf("client canceled request\n")
			return nil, status.Error(codes.Canceled, "request canceled")
		}
		time.Sleep(time.Second)
	}

	firstname := req.GetGreeting().GetFirstName()
	result := "Hello " + firstname
	res := &greetwithdeadlinepb.GreetWithDeadlineResponse{
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
	greetwithdeadlinepb.RegisterGreetServiceServer(server, &serviceServer{})

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Serve() failed: %v\n", err)
	}
}
