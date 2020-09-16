package main

import (
	"context"
	"fmt"
	"io"
	"log"

	blogpb "github.com/MickiFoerster/grpc/blog"
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
	clientConn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Dial() failed: %v\n", err)
	}
	defer clientConn.Close()

	serviceClient := blogpb.NewBlogServiceClient(clientConn)
	doUnary(serviceClient)
}

func doUnary(serviceClient blogpb.BlogServiceClient) {
	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "John Doo",
			Content:  "strange things",
			Title:    "Content of the first blog",
		},
	}
	res, err := serviceClient.CreateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("CreateBlog() failed: %v\n", err)
	}
	fmt.Printf("server response: %v\n", res.GetBlog())
	createdBlogId := res.GetBlog().GetId()

	fmt.Println("Reading blog entries")

	// read arbitrary blog -> error
	_, err = serviceClient.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "1234gasdf"})
	if err != nil {
		fmt.Printf("could not read document: %v\n", err)
	}

	// read blog from above
	match, err := serviceClient.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: createdBlogId})
	if err != nil {
		fmt.Printf("could not read document: %v\n", err)
	} else {
		fmt.Printf("blog found: %v\n", match)
	}

	// update blog
	newBlog := &blogpb.Blog{
		Id:       createdBlogId,
		AuthorId: "Spiderman",
		Content:  "Comics",
		Title:    "Exciting stories",
	}
	update_result, update_err := serviceClient.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if update_err != nil {
		log.Printf("error: update failed: %v\n", update_err)
	} else {
		fmt.Printf("Blog updated: %v\n", update_result)
	}

	// delete blog
	delete_result, delete_err := serviceClient.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: createdBlogId})
	if delete_err != nil {
		log.Printf("error: delete failed: %v\n", delete_err)
	} else {
		fmt.Printf("Blog deleted: %v\n", delete_result)
	}

	// list blog
	stream, err := serviceClient.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error: RPC ListBlog failed: %v\n", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error: while receiving stream: %v\n", err)
		}
		fmt.Println(res.GetBlog())
	}
}
