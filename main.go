package main

import (
	"log"
	"net"

	pb "github.com/claudiocleberson/shippy-service-consignment/proto/consignment"
	"github.com/claudiocleberson/shippy-service-consignment/repository"
	"github.com/claudiocleberson/shippy-service-consignment/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

func main() {
	repo := repository.NewRepository()
	srvConsignment := services.NewService(repo)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	pb.RegisterShippingServiceServer(s, srvConsignment)

	reflection.Register(s)

	log.Println("Running server on port:", port)
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
