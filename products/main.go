package main

import (
	"log"
	"net"
	"strings"

	// Import the generated protobuf code
	pb "github.com/barnie/kibubu/products/genproto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	port = ":50051"
)

// IRepository is the interface for repository
type IRepository interface {
	Create(*pb.Product) (*pb.Product, error)
	GetAll() ([]*pb.Product, error)
	GetProduct(id string) (*pb.Product, error)
}

// Repository - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later on.
type Repository struct {
	products []*pb.Product
}

// Create : Method create for database
func (repo *Repository) Create(product *pb.Product) (*pb.Product, error) {
	updated := append(repo.products, product)
	repo.products = updated
	return product, nil
}

// GetAll : Method create for database
func (repo *Repository) GetAll() ([]*pb.Product, error) {
	return repo.products, nil
}

// GetProduct : Method create for database
func (repo *Repository) GetProduct(id string) (*pb.Product, error) {
	var found *pb.Product
	for i := 0; i < len(repo.products); i++ {
		if id == repo.products[i].Id {
			found = repo.products[i]
		}
	}
	if found == nil {
		return nil, status.Errorf(codes.NotFound, "no product with ID %s", id)
	}
	return found, nil

}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
type service struct {
	repo IRepository
}

// CreateProduct - we created just one method on our service,
func (s *service) CreateProduct(ctx context.Context, req *pb.Product) (*pb.CreateProductResponse, error) {

	// Save our product
	product, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	// Return matching the `CreateProductResponse` message we created in our
	// protobuf definition.
	return &pb.CreateProductResponse{Created: true, Product: product}, nil
}

// ListProducts - we created just one method on our service
func (s *service) ListProducts(ctx context.Context, req *pb.Empty) (*pb.ListProductsResponse, error) {

	// Save our product
	products, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	// Return matching the `ListProductsResponse` message we created in our
	// protobuf definition.
	return &pb.ListProductsResponse{Products: products}, nil
}

// GetProduct - we created just one method on our service
func (s *service) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {

	// Save our product
	product, err := s.repo.GetProduct(req.Id)
	if err != nil {
		return nil, err
	}

	// Return matching the `Products` message we created in our
	// protobuf definition.
	return product, nil
}

// CreateProduct - we created just one method on our service
func (s *service) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {

	// Save our product
	var ps []*pb.Product
	var results, _ = s.repo.GetAll()
	for _, p := range results {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(req.Query)) ||
			strings.Contains(strings.ToLower(p.Description), strings.ToLower(req.Query)) {
			ps = append(ps, p)
		}
	}

	// Return matching the `SearchProductsResponse` message we created in our
	// protobuf definition.
	return &pb.SearchProductsResponse{Results: ps}, nil
}

func main() {

	repo := &Repository{}

	// Set-up our gRPC server.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	pb.RegisterProductServiceServer(s, &service{repo})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
