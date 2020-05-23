package main

import (
	pb "github.com/claudiocleberson/shippy-service-consignment/proto/consignment"
	"github.com/claudiocleberson/shippy-service-consignment/repository"
	"github.com/claudiocleberson/shippy-service-consignment/services"
	"github.com/micro/go-micro"
)

const (
	port = ":50051"
)

func main() {
	repo := repository.NewRepository()
	srvConsignment := services.NewService(repo)

	srv := micro.NewService(
		micro.Name("shippy.service.consignment"),
	)
	srv.Init()

	pb.RegisterShippingServiceHandler(srv.Server(), srvConsignment)

	if err := srv.Run(); err != nil {
		panic(err)
	}

	// lis, err := net.Listen("tcp", port)
	// if err != nil {
	// 	panic(err)
	// }

	// s := grpc.NewServer()

	// pb.RegisterShippingServiceServer(s, srvConsignment)

	// reflection.Register(s)

	// log.Println("Running server on port:", port)
	// if err := s.Serve(lis); err != nil {
	// 	panic(err)
	// }
}
