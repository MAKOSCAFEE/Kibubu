package main

import (
	"context"
	"log"

	pb "github.com/barnie/kibubu/products/genproto"
	"google.golang.org/grpc"
)

const (
	address         = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewProductServiceClient(conn)

	getAll, err := client.ListProducts(context.Background(), &pb.Empty{})

	if err != nil {
		log.Fatalf("Could not list consignments: %v", err)
	}
	for _, v := range getAll.Products {
		log.Println(v)
	}

}
