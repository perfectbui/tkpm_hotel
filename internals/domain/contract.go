package domain

import (
	"TKPM/internals/models"
	"TKPM/internals/repository"
	"context"
)

type Contract interface {
	Create(ctx context.Context, Contract *models.Contract) (*models.Contract, error)
	Update(ctx context.Context, Contract *models.Contract) (*models.Contract, error)
	GetContractList(ctx context.Context, offset, limit int64) ([]*models.Contract, error)
	GetContractByUserId(ctx context.Context, userId string) (*models.Contract, error)
}

type contractDomain struct {
	contractRepository repository.Contract
}

func NewContractDomain(contractRepository repository.Contract) Contract {
	return &contractDomain{
		contractRepository: contractRepository,
	}
}

func (d *contractDomain) Create(ctx context.Context, contract *models.Contract) (*models.Contract, error) {
	res, err := d.contractRepository.Create(ctx, contract)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *contractDomain) Update(ctx context.Context, contract *models.Contract) (*models.Contract, error) {
	res, err := d.contractRepository.Update(ctx, contract)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *contractDomain) GetContractList(ctx context.Context, offset, limit int64) ([]*models.Contract, error) {
	res, err := d.contractRepository.GetContractList(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *contractDomain) GetContractByUserId(ctx context.Context, userId string) (*models.Contract, error) {
	res, err := d.contractRepository.GetContractByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}
