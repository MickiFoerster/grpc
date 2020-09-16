package main

import (
	"fmt"
	"io"
	"log"
	"net"

	greetbidirectionalstreampb "github.com/MickiFoerster/grpc/greet-bidirectional-stream"
	grpc "google.golang.org/grpc"
)

type serviceServer struct {
}

func (*serviceServer) GreetServerStream(stream greetbidirectionalstreampb.GreetService_GreetServerStreamServer) error {
	log.Printf("RPC call received\n")

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		lastName := req.GetGreeting().GetLastName()
		result := "Hello " + firstName + " " + lastName + "!\n"

		log.Printf("send message to client: %v\n", result)
		err = stream.Send(&greetbidirectionalstreampb.GreetResponse{
			Result: result,
		})
		if err != nil {
			log.Printf("error while sending response to client: %v\n", err)
		}
	}
}

func main() {
	fmt.Println("Go server starts ...")

	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("listen failed: %v\n", err)
	}

	srv := grpc.NewServer()
	greetbidirectionalstreampb.RegisterGreetServiceServer(srv, &serviceServer{})

	if err := srv.Serve(listener); err != nil {
		log.Fatalf("Serve() failed: %v\n", err)
	}
	log.Printf("server finished\n")
}
