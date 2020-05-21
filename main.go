package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/claudiocleberson/shippy-service-consignment/shippy-service-consignment/proto/consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

type Repository struct {
	mu           sync.RWMutex
	consignments []*pb.Consignment
}

func (repo *Repository) Create(cons *pb.Consignment) (*pb.Consignment, error) {
	repo.mu.Lock()
	update := append(repo.consignments, cons)
	repo.consignments = update
	repo.mu.Unlock()
	return cons, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {

	return repo.consignments
}

type service struct {
	repo repository
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {

	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	resp := pb.Response{
		Created:     true,
		Consignment: consignment,
	}

	return &resp, nil

}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	consignments := s.repo.GetAll()

	response := &pb.Response{
		ListConsignments: consignments,
	}
	return response, nil
}

func main() {
	repo := &Repository{}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	pb.RegisterShippingServiceServer(s, &service{
		repo: repo})

	reflection.Register(s)

	log.Println("Running server on port:", port)
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
