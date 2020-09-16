package main

import (
	"context"
	"fmt"
	"log"

	"github.com/MickiFoerster/grpc/greet/go"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("client starts ...")

	clientConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial() failed: %v\n", err)
	}
	defer clientConn.Close()

	serviceClient := greetpb.NewGreetServiceClient(clientConn)
	//fmt.Printf("client created: %v\n", serviceClient)
	doUnary(serviceClient)
}

func doUnary(serviceClient greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName:  "Doo",
		},
	}
	res, err := serviceClient.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Greet() failed: %v\n", err)
	}
	fmt.Printf("server response: %q\n", res.Result)
}
