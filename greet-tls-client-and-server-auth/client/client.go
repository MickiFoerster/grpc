package main

import (
	"context"
	"fmt"
	"log"

	greettlspb "github.com/MickiFoerster/grpc/greet-tls-client-and-server-auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	fmt.Println("client starts ...")

	certFile := "ca.crt" // certificate of Certificate Authority
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		log.Fatalf("error while creating TLS client: %v\n", err)
	}

	opts := grpc.WithTransportCredentials(creds)
	clientConn, err := grpc.Dial("localhost:55551", opts)
	if err != nil {
		log.Fatalf("Dial() failed: %v\n", err)
	}
	defer clientConn.Close()

	serviceClient := greettlspb.NewGreetServiceClient(clientConn)
	//fmt.Printf("client created: %v\n", serviceClient)
	doUnary(serviceClient)
}

func doUnary(serviceClient greettlspb.GreetServiceClient) {
	req := &greettlspb.GreetRequest{
		Greeting: &greettlspb.Greeting{
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
