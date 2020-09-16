package main

import (
	"context"
	"fmt"
	"log"
	"time"

	greetwithdeadlinepb "github.com/MickiFoerster/grpc/greet-with-deadline"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("client starts ...")

	clientConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial() failed: %v\n", err)
	}
	defer clientConn.Close()

	serviceClient := greetwithdeadlinepb.NewGreetServiceClient(clientConn)
	//fmt.Printf("client created: %v\n", serviceClient)
	doUnary(serviceClient)
}

func doUnary(serviceClient greetwithdeadlinepb.GreetServiceClient) {
	req := &greetwithdeadlinepb.GreetWithDeadlineRequest{
		Greeting: &greetwithdeadlinepb.Greeting{
			FirstName: "John",
			LastName:  "Doo",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := serviceClient.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Printf("RPC took too long, canceled\n")
			} else {
				fmt.Printf("unexpected error: %v\n", statusErr)
			}
		} else {
			log.Fatalf("Greet() failed: %v\n", err)
		}
	} else {
		fmt.Printf("server response: %q\n", res.Result)
	}
}
