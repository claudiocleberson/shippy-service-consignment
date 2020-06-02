package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/claudiocleberson/shippy-service-consignment/datastore"
	"github.com/claudiocleberson/shippy-service-consignment/handlers"
	pb "github.com/claudiocleberson/shippy-service-consignment/proto/consignment"
	"github.com/claudiocleberson/shippy-service-consignment/repository"
	userService "github.com/claudiocleberson/shippy-service-users/proto/users"
	vesselProto "github.com/claudiocleberson/shippy-service-vessel/proto/vessel"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
)

const (
	port        = ":50051"
	dbHost      = "DB_HOST"
	defaultHost = "mongodb://localhost:27017"
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
		micro.Version("latest"),

		//Out auth middleware
		micro.WrapHandler(AuthWrapper),
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

func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.New("Missing token on the request!")
		}

		//Note this now uppercase
		token := meta["Token"]
		log.Println("Authenticatin with token: ", token)

		//Validate token
		authClient := userService.NewUserServiceClient("shippy.service.users", client.DefaultClient)
		_, err := authClient.ValidateToken(context.Background(), &userService.Token{
			Token: token,
		})

		if err != nil {
			return err
		}

		err = fn(ctx, req, resp)
		return err

	}
}
