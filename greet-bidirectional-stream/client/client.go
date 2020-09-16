package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	greetbidirectionalstreampb "github.com/MickiFoerster/grpc/greet-bidirectional-stream"
	grpc "google.golang.org/grpc"
)

func main() {
	fmt.Println("client starts ...")

	clientConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial() failed: %v\n", err)
	}
	defer clientConn.Close()

	serviceClient := greetbidirectionalstreampb.NewGreetServiceClient(clientConn)
	doBiDirectionalStreaming(serviceClient)
}

func doBiDirectionalStreaming(client greetbidirectionalstreampb.GreetServiceClient) {
	fmt.Println("Starting to do a bidirectional streaming RPC ...")

	stream, err := client.GreetServerStream(context.Background())
	if err != nil {
		log.Fatalf("error while calling RPC: %v", err)
	}

	done := make(chan struct{})
	// Producer routine
	go func() {
		type name struct {
			firstname string
			lastname  string
		}
		names := []name{
			name{firstname: "Barack", lastname: "Obama"},
			name{firstname: "James", lastname: "Cameroon"},
			name{firstname: "John", lastname: "Doo"},
		}
		for _, n := range names {
			fmt.Printf("Sending message: %v", n)
			stream.Send(&greetbidirectionalstreampb.GreetRequest{
				Greeting: &greetbidirectionalstreampb.Greeting{
					FirstName: n.firstname,
					LastName:  n.lastname,
				},
			})
			time.Sleep(time.Second)
		}
		err := stream.CloseSend()
		if err != nil {
			log.Printf("error while closing stream\n", err)
		}
	}()

	// Consumer of server answers
	go func() {
		for {
			res, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					close(done)
					return
				}
				log.Printf("error while receiving server response: %v\n", err)
				continue
			}
			fmt.Printf("server response: %v\n", res.GetResult())
		}
	}()

	<-done
	fmt.Println("Bi directional stream closed")
}
