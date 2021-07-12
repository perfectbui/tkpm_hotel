package repository

import (
	"TKPM/common/enums"
	"TKPM/internals/models"
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Contract interface {
	Create(ctx context.Context, Contract *models.Contract) (*models.Contract, error)
	Update(ctx context.Context, Contract *models.Contract) (*models.Contract, error)
	GetContractList(ctx context.Context, offset, limit int64) ([]*models.Contract, error)
	GetContractByUserId(ctx context.Context, userId string) (*models.Contract, error)
}

type contractRepository struct {
	coll *mongo.Collection
}

func NewContractRepository(coll *mongo.Collection) Contract {
	return &contractRepository{
		coll: coll,
	}
}

func (r *contractRepository) Create(ctx context.Context, contract *models.Contract) (*models.Contract, error) {
	contract.ContractID = uuid.New().String()
	t := time.Now()
	contract.CreatedAt = &t
	contract.Status = enums.Processing.String()
	_, err := r.coll.InsertOne(ctx, contract)
	return contract, err
}

func (r *contractRepository) Update(ctx context.Context, contract *models.Contract) (*models.Contract, error) {
	after := options.After
	res := r.coll.FindOneAndUpdate(ctx, bson.M{"contract_id": contract.ContractID}, bson.M{
		"$set": contract,
	}, &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}

	var updatedContract models.Contract
	err := res.Decode(&updatedContract)
	if err != nil {
		return nil, err
	}

	return &updatedContract, nil
}

func (r *contractRepository) GetContractList(ctx context.Context, offset, limit int64) ([]*models.Contract, error) {
	result := []*models.Contract{}

	opt := options.FindOptions{
		Skip:  &offset,
		Limit: &limit,
	}
	cur, err := r.coll.Find(ctx, &models.Contract{}, &opt)

	if err != nil {
		return nil, err
	}

	if cur.Err() != nil {
		return nil, cur.Err()
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *contractRepository) GetContractByUserId(ctx context.Context, userId string) (*models.Contract, error) {
	var contract models.Contract
	err := r.coll.FindOne(ctx, bson.M{"user_id": userId}).Decode(&contract)
	if err != nil {
		return nil, err
	}

	return &contract, nil
}
