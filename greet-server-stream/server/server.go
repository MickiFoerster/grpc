package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	greetserverstreampb "github.com/MickiFoerster/grpc/greet-stream"
	grpc "google.golang.org/grpc"
)

type serviceServer struct {
}

func (*serviceServer) GreetServerStream(req *greetserverstreampb.GreetRequest,
	stream greetserverstreampb.GreetService_GreetServerStreamServer) error {
	fmt.Printf("GreetServerStream() was invoked with %v\n", req)
	firstname := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstname + " number " + strconv.Itoa(i)
		res := &greetserverstreampb.GreetResponse{
			Result: result,
		}
		stream.Send(res)
		log.Printf("Server now starts to prepare next response ... ")
		time.Sleep(time.Second)
		log.Println("DONE")
	}
	log.Printf("Closing stream\n")
	return nil
}

func main() {
	fmt.Println("Go server starts ...")

	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("listen failed: %v\n", err)
	}

	srv := grpc.NewServer()
	greetserverstreampb.RegisterGreetServiceServer(srv, &serviceServer{})

	if err := srv.Serve(listener); err != nil {
		log.Fatalf("Serve() failed: %v\n", err)
	}
}
