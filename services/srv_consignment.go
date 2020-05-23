package services

import (
	"context"
	"fmt"

	pb "github.com/claudiocleberson/shippy-service-consignment/proto/consignment"
	repo "github.com/claudiocleberson/shippy-service-consignment/repository"
	vesselProto "github.com/claudiocleberson/shippy-service-vessel/proto/vessel"
)

type ServiceConsignmentInterface interface {
	CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error
	GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error
}

type serviceConsignment struct {
	repo         repo.ConsignmentRepository
	vesselClient vesselProto.VesselServiceClient
}

func NewService(repo repo.ConsignmentRepository, vesselCli vesselProto.VesselServiceClient) ServiceConsignmentInterface {
	return &serviceConsignment{
		repo:         repo,
		vesselClient: vesselCli,
	}
}

func (s *serviceConsignment) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

	//Here we call a client instance of our vessel service with our consignment weight,
	//and the amount of containers as the capacity value
	vesselResponse, err := s.vesselClient.FindAvailable(ctx, &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})

	if err != nil {
		return nil
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

func (s *serviceConsignment) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {

	consignments := s.repo.GetAll()

	res.ListConsignments = consignments

	return nil
}
