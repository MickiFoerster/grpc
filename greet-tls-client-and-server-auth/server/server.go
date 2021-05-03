// TLS server for grpc
// Details under https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-auth-support.md
//               https://grpc.io/docs/guides/auth.html
package main

import (
	"context"
	"fmt"
	"log"
	"net"

	greettlspb "github.com/MickiFoerster/grpc/greet-tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type serviceServer struct {
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

	listener, err := net.Listen("tcp", "0.0.0.0:50051")
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
