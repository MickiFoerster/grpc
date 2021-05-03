// TLS server for grpc
// Details under https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-auth-support.md
//               https://grpc.io/docs/guides/auth.html
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"

	greettlspb "github.com/MickiFoerster/grpc/greet-tls-client-and-server-auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type serviceServer struct {
	greettlspb.UnimplementedGreetServiceServer
}

func (*serviceServer) Greet(ctx context.Context, req *greettlspb.GreetRequest) (*greettlspb.GreetResponse, error) {
	fmt.Printf("Greet() was invoked with %v\n", req)
	firstname := req.GetGreeting().GetFirstName()
	result := "Hello " + firstname
	res := &greettlspb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Go server starts ...")

	listener, err := net.Listen("tcp", "0.0.0.0:55551")
	if err != nil {
		log.Fatalf("listen failed: %v\n", err)
	}

	certFile := "server.crt"
	keyFile := "server.pem"
	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	if sslErr != nil {
		log.Fatalf("error while loading certificates: %v\n", sslErr)
	}

	opts := grpc.Creds(creds)

	server := grpc.NewServer(opts)
	greettlspb.RegisterGreetServiceServer(server, &serviceServer{})

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Serve() failed: %v\n", err)
	}
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(
		"server.pem",
		"server.key")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}
