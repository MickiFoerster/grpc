package main

import (
	"context"
	"fmt"
	"io"
	"log"

	greetserverstreampb "github.com/MickiFoerster/grpc/greet-stream"
	grpc "google.golang.org/grpc"
)

func main() {
	fmt.Println("client starts ...")

	clientConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial() failed: %v\n", err)
	}
	defer clientConn.Close()

	serviceClient := greetserverstreampb.NewGreetServiceClient(clientConn)
	doServerStreaming(serviceClient)
}

func doServerStreaming(client greetserverstreampb.GreetServiceClient) {
	fmt.Println("Starting to do a server streaming RPC ...")
	req := &greetserverstreampb.GreetRequest{
		Greeting: &greetserverstreampb.Greeting{
			FirstName: "hello",
			LastName:  "world",
		},
	}
	stream, err := client.GreetServerStream(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling RPC: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// end of stream reached
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		fmt.Printf("Response from GreetManyTimes: %v\n", msg.GetResult())
	}
}
