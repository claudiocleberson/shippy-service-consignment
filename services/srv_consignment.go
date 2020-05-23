package services

import (
	"context"

	pb "github.com/claudiocleberson/shippy-service-consignment/proto/consignment"
	repo "github.com/claudiocleberson/shippy-service-consignment/repository"
)

type ServiceConsignmentInterface interface {
	CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error
	GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error
}

type serviceConsignment struct {
	repo repo.ConsignmentRepository
}

func NewService(repo repo.ConsignmentRepository) ServiceConsignmentInterface {
	return &serviceConsignment{
		repo: repo,
	}
}

func (s *serviceConsignment) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

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
