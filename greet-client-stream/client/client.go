package main

import (
	"context"
	"fmt"
	"log"
	"time"

	greetclientstreampb "github.com/MickiFoerster/grpc/greet-client-stream"
	grpc "google.golang.org/grpc"
)

func main() {
	fmt.Println("client starts ...")

	clientConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial() failed: %v\n", err)
	}
	defer clientConn.Close()

	serviceClient := greetclientstreampb.NewGreetServiceClient(clientConn)
	doClientStreaming(serviceClient)
}

func doClientStreaming(client greetclientstreampb.GreetServiceClient) {
	fmt.Println("Starting to do a stream from client to server ...")
	stream, err := client.GreetClientStream(context.Background())
	if err != nil {
		log.Fatalf("error while streaming from client to server: %v", err)
	}

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
		log.Printf("sending %v ...\n", n)
		stream.Send(&greetclientstreampb.GreetRequest{
			Greeting: &greetclientstreampb.Greeting{
				FirstName: n.firstname,
				LastName:  n.lastname,
			},
		})
		time.Sleep(time.Second)
	}

	log.Printf("client ends stream\n")
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while closing stream: %v", err)
	}
	log.Printf("client stream has been ended. Server answer: %v", res.GetResult())
}
