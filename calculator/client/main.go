package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	calculatorpb "github.com/MickiFoerster/grpc/calculator"
	grpc "google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	client := calculatorpb.NewCalculatorServiceClient(cc)
	//doUnaryRequestSum(client)
	doUnaryRequestSquareRoot(client, 64)
	doUnaryRequestSquareRoot(client, -64)
}

func doUnaryRequestSquareRoot(client calculatorpb.CalculatorServiceClient, number float64) {
	req := &calculatorpb.SquareRootRequest{
		Number: number,
	}

	fmt.Printf("Send Squareroot RPC with %v as input\n", number)
	res, err := client.SquareRoot(context.Background(), req)
	if err != nil {
		responseStatus, ok := status.FromError(err)
		if ok {
			// user defined err
			fmt.Printf("%s: %v\n", responseStatus.Message(), responseStatus.Code())
			if responseStatus.Code() == codes.InvalidArgument {
				fmt.Printf("error: You sent a negative number: %v\n", number)
			}
		} else {
			fmt.Printf("grpc error while sending squareroot RPC: %v\n", err)
		}
		return
	}
	fmt.Printf("Square root of %v is %v.\n", number, res.GetNumberRoot())
}

func doUnaryRequestSum(client calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SumRequest{
		FirstNumber:  3,
		SecondNumber: 5,
	}
	fmt.Println("Trigger a unary RPC call")
	res, err := client.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling RPC: %v", err)
	}
	log.Printf("RPC response: %v", res.GetSumResult())
}
