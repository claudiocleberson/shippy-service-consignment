package repository

import (
	"sync"

	pb "github.com/claudiocleberson/shippy-service-consignment/proto/consignment"
)

type ConsignmentRepository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

type consignmentRepository struct {
	mu           sync.RWMutex
	consignments []*pb.Consignment
}

func NewRepository() ConsignmentRepository {
	return &consignmentRepository{}
}

func (repo *consignmentRepository) Create(cons *pb.Consignment) (*pb.Consignment, error) {
	repo.mu.Lock()
	update := append(repo.consignments, cons)
	repo.consignments = update
	repo.mu.Unlock()
	return cons, nil
}

func (repo *consignmentRepository) GetAll() []*pb.Consignment {

	return repo.consignments
}
