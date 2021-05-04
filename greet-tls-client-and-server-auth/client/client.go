package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	greettlspb "github.com/MickiFoerster/grpc/greet-tls-client-and-server-auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	fmt.Println("client starts ...")

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatalln("Cannot load TLS credentials:", err)
	}

	opts := grpc.WithTransportCredentials(tlsCredentials)
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

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Load server's certificate and private key
	clientCert, err := tls.LoadX509KeyPair("client.crt", "client.pem")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{clientCert},
	}

	return credentials.NewTLS(config), nil
}
