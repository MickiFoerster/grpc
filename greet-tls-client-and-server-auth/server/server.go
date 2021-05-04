// TLS server for grpc
// Details under https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-auth-support.md
//               https://grpc.io/docs/guides/auth.html
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
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

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatalln("Cannot load TLS credentials", err)
	}

	opts := grpc.Creds(tlsCredentials)
	server := grpc.NewServer(opts)
	greettlspb.RegisterGreetServiceServer(server, &serviceServer{})

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Serve() failed: %v\n", err)
	}
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed client's certificate
	pemClientCA, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}

	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("server.crt", "server.pem")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}
