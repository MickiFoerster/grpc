package main

import (
	"fmt"
	"io"
	"log"
	"net"

	greetclientstreampb "github.com/MickiFoerster/grpc/greet-client-stream"
	grpc "google.golang.org/grpc"
)

type serviceServer struct {
}

func (*serviceServer) GreetClientStream(stream greetclientstreampb.GreetService_GreetClientStreamServer) error {
	log.Printf("received RPC call, start handling request ...\n")
	result := ""
	for {
		req, err := stream.Recv()
		if req == nil && err == nil {
			log.Printf("received empty request, continue ...")
			continue
		}

		if err != nil {
			if err == io.EOF {
				log.Printf("client has finished the stream, req=%v, err=%v\n", req, err)
				stream.SendAndClose(&greetclientstreampb.GreetResponse{
					Result: result,
				})
				break
			}
			log.Fatalf("error while reading client stream: %v\n", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		log.Printf("received %v from stream\n", firstName)
		result += "Hello " + firstName + "!\n"
	}
	log.Printf("server ends loop where receiving stream\n")
	return nil
}

func main() {
	fmt.Println("Go client starts ...")

	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("listen failed: %v\n", err)
	}

	srv := grpc.NewServer()
	greetclientstreampb.RegisterGreetServiceServer(srv, &serviceServer{})

	if err := srv.Serve(listener); err != nil {
		log.Fatalf("Serve() failed: %v\n", err)
	}
}
