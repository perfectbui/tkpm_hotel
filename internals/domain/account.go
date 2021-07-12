package domain

import (
	"TKPM/common/errors"
	"TKPM/internals/models"
	"TKPM/internals/repository"
	"context"
)

type Account interface {
	SignUp(ctx context.Context, req *models.Account) (*models.Account, error)
	SignIn(ctx context.Context, req *models.Account) (*models.Account, error)
	CheckAuth(ctx context.Context, accountId string) (*models.Account, error)
}

type accountDomain struct {
	accountRepository repository.Account
}

func NewAccountDomain(accountRepository repository.Account) Account {
	return &accountDomain{
		accountRepository: accountRepository,
	}
}

func (d *accountDomain) SignUp(ctx context.Context, req *models.Account) (*models.Account, error) {
	res, err := d.accountRepository.FindByEmail(ctx, req.Email)
	if res != nil {
		return nil, err
	}

	res, err = d.accountRepository.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	
	return &models.Account{
		AccountID: res.AccountID,
	}, nil
}

func (d *accountDomain) SignIn(ctx context.Context, req *models.Account) (*models.Account, error) {
	res, err := d.accountRepository.FindByEmail(ctx, req.Email)
	if res == nil {
		return nil, err
	}

	if !res.IsCorrectPassword(req.Password) {
		return nil, errors.ErrPasswordIsNotCorrect
	}

	return res, nil
}

func (d *accountDomain) CheckAuth(ctx context.Context, accountId string) (*models.Account, error) {
	res, err := d.accountRepository.FindByAccountId(ctx, accountId)
	if err != nil {
		return nil, err
	}

	return res, nil
}
