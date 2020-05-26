package main

import (
	"os"

	"github.com/claudiocleberson/shippy-service-consignment/datastore"
	"github.com/claudiocleberson/shippy-service-consignment/handlers"
	pb "github.com/claudiocleberson/shippy-service-consignment/proto/consignment"
	"github.com/claudiocleberson/shippy-service-consignment/repository"
	vesselProto "github.com/claudiocleberson/shippy-service-vessel/proto/vessel"
	"github.com/micro/go-micro"
)

const (
	port        = ":50051"
	dbHost      = "DB_HOST"
	defaultHost = "mongodb://datastore:27017"
)

func main() {

	mongoUri := os.Getenv(dbHost)
	if mongoUri == "" {
		mongoUri = defaultHost
	}

	dbClient := datastore.NewMongoClient(mongoUri)

	repo := repository.NewRepository(dbClient)

	srv := micro.NewService(
		micro.Name("shippy.service.consignment"),
	)
	srv.Init()

	//Declare service clientes dependecies
	vesselClient := vesselProto.NewVesselServiceClient("shippy.service.vessel", srv.Client())

	srvConsignment := handlers.NewConsignmentHandler(repo, vesselClient)

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
