// TLS server for grpc
// Details under https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-auth-support.md
//               https://grpc.io/docs/guides/auth.html
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	blogpb "github.com/MickiFoerster/grpc/blog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

var collection *mongo.Collection

type serviceServer struct {
}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func (*serviceServer) ListBlog(req *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {
	log.Println("Received RPC Listblog")

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Find failed: %v\n", err))
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		data := &blogItem{}
		err := cursor.Decode(data)
		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("Decode failed: %v\n", err))
		}
		// send data to client
		stream.Send(&blogpb.ListBlogResponse{Blog: &blogpb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorID,
			Content:  data.Content,
			Title:    data.Title,
		}})
	}

	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Cursor error: %v\n", err))
	}

	return nil
}

func (*serviceServer) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	log.Println("Received RPC DeleteBlog")
	oid, err := primitive.ObjectIDFromHex(req.GetBlogId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid ID provided"))
	}

	filter := bson.D{{"_id", oid}}
	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("cannot delete blog ID %v: %v\n", req.GetBlogId()), err)
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("cannot find blog ID %v: %v\n", req.GetBlogId(), err))
	}

	return &blogpb.DeleteBlogResponse{BlogId: req.GetBlogId()}, nil
}

func (*serviceServer) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	log.Println("Received RPC UpdateBlog")
	blog := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid ID provided"))
	}

	data := &blogItem{}
	filter := bson.D{{"_id", oid}}
	res := collection.FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("cannot find ID %q in MongoDB database\n", blog.GetId()))
	}

	data.AuthorID = blog.GetAuthorId()
	data.Content = blog.GetContent()
	data.Title = blog.GetTitle()

	_, err = collection.ReplaceOne(ctx, filter, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("cannot update object with ID %q in MongoDB database: %v\n", oid, err))
	}

	return &blogpb.UpdateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorID,
			Content:  data.Content,
			Title:    data.Title,
		},
	}, nil
}

func (*serviceServer) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	fmt.Println("request ReadBlog received")

	blogId := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid ID provided"))
	}

	data := &blogItem{}
	filter := bson.D{{"_id", oid}}
	res := collection.FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("cannot find ID %q in MongoDB database\n", blogId))
	}

	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorID,
			Content:  data.Content,
			Title:    data.Title,
		},
	}, nil
}

func (*serviceServer) CreateBlog(ctx context.Context,
	req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {

	blog := req.GetBlog()

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetTitle(),
		Title:    blog.GetContent(),
	}

	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v\n", err)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			"cannot convert to OID")
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Content:  blog.GetTitle(),
			Title:    blog.GetContent(),
		},
	}, nil

}

//func (*serviceServer) Greet(ctx context.Context, req *blogpb.GreetRequest) (*blogpb.GreetResponse, error) {
//	fmt.Printf("Greet() was invoked with %v\n", req)
//	firstname := req.GetGreeting().GetFirstName()
//	result := "Hello " + firstname
//	res := &blogpb.GreetResponse{
//		Result: result,
//	}
//	return res, nil
//}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// start signal handling
	ch := make(chan os.Signal, 1)
	var listener net.Listener
	var server *grpc.Server
	var client *mongo.Client
	var ctx context.Context
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		if server != nil {
			log.Println("signal handler stops server")
			server.Stop()
		}
		if listener != nil {
			log.Println("signal handler stops listener")
			listener.Close()
		}
		if client != nil {
			log.Println("stop MongoDB client")
			client.Disconnect(ctx)
		}
	}()

	// open MongoDB connection
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("error while creating mongodb client: %v\n", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("error while creating mongodb client: %v\n", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("error: could not connect to MongoDB server. Is MongoDB server running?")
	}

	// create/open MongoDB database
	collection = client.Database("mydb").Collection("blog")

	l, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("listen failed: %v\n", err)
	}
	listener = l

	certFile := "server.crt"
	keyFile := "server.pem"
	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	if sslErr != nil {
		log.Fatalf("error while loading certificates: %v\n", sslErr)
	}

	opts := grpc.Creds(creds)

	server = grpc.NewServer(opts)
	blogpb.RegisterBlogServiceServer(server, &serviceServer{})

	fmt.Println("Blog service started ...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Serve() failed: %v\n", err)
	}
}
