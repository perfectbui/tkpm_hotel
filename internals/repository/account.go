package repository

import (
	"TKPM/internals/models"
	"context"

	"TKPM/common/enums"
	"TKPM/common/errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Account interface {
	Create(ctx context.Context, account *models.Account) (*models.Account, error)
	FindByEmail(ctx context.Context, email string) (*models.Account, error)
	FindByAccountId(ctx context.Context, accountId string) (*models.Account, error)
}

type accountRepository struct {
	coll *mongo.Collection
}

func NewAccountRepository(coll *mongo.Collection) Account {
	return &accountRepository{
		coll: coll,
	}
}

func (r *accountRepository) Create(ctx context.Context, account *models.Account) (*models.Account, error) {
	err := account.HashPassword()
	account.AccountID = uuid.New().String()
	account.Role = enums.User.String()
	if err != nil {
		return nil, err
	}

	_, err = r.coll.InsertOne(ctx, account)
	return account, err
}

func (r *accountRepository) FindByEmail(ctx context.Context, email string) (*models.Account, error) {
	result := models.Account{}
	err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, errors.ErrEmailNotFound
	}

	if err != nil {
		return nil, err
	}
	return &result, errors.ErrEmailExisted
}

func (r *accountRepository) FindByAccountId(ctx context.Context, accountId string) (*models.Account, error) {
	result := models.Account{}
	err := r.coll.FindOne(ctx, bson.M{"account_id": accountId}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
