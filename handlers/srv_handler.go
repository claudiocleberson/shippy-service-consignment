package handlers

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/claudiocleberson/shippy-service-consignment/proto/consignment"
	repo "github.com/claudiocleberson/shippy-service-consignment/repository"
	vesselProto "github.com/claudiocleberson/shippy-service-vessel/proto/vessel"
)

type HandlerConsignmentInterface interface {
	CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error
	GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error
}

type handlerConsignment struct {
	repo         repo.ConsignmentRepository
	vesselClient vesselProto.VesselServiceClient
}

func NewConsignmentHandler(repo repo.ConsignmentRepository, vesselCli vesselProto.VesselServiceClient) HandlerConsignmentInterface {
	return &handlerConsignment{
		repo:         repo,
		vesselClient: vesselCli,
	}
}

func (s *handlerConsignment) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

	//Here we call a client instance of our vessel service with our consignment weight,
	//and the amount of containers as the capacity value
	vesselResponse, err := s.vesselClient.FindAvailable(ctx, &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})

	if err != nil {
		return errors.New("no vessel available")
	}

	fmt.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)

	//Set the vesselId as the vessel we got back from out vessel service.
	req.VesselId = vesselResponse.Vessel.Id

	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	res.Created = true
	res.Consignment = consignment

	return nil

}

func (s *handlerConsignment) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {

	consignments := s.repo.GetAll(ctx)

	res.ListConsignments = consignments

	return nil
}
