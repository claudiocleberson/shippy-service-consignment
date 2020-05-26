package repository

import (
	"context"

	"github.com/claudiocleberson/shippy-service-consignment/datastore"
	"github.com/claudiocleberson/shippy-service-consignment/models"
	pb "github.com/claudiocleberson/shippy-service-consignment/proto/consignment"
)

type ConsignmentRepository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

type consignmentRepository struct {
	mongoClient datastore.MongoClient
}

func NewRepository(datastore datastore.MongoClient) ConsignmentRepository {
	return &consignmentRepository{
		mongoClient: datastore,
	}
}

func (repo *consignmentRepository) Create(cons *pb.Consignment) (*pb.Consignment, error) {

	ctx := context.Background()

	err := repo.mongoClient.Create(ctx, models.MarshalConsignment(cons))
	if err != nil {
		return nil, err
	}

	return cons, nil
}

func (repo *consignmentRepository) GetAll() []*pb.Consignment {

	ctx := context.Background()
	result, err := repo.mongoClient.GetAll(ctx)
	if err != nil {
		return nil
	}

	return models.UnmarshalConsignmentCollection(result)
}
