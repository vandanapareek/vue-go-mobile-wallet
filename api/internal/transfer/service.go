package transfer

import (
	"context"
	"database/sql"
	"fmt"
	"thunes-api/errors"
)

type Service interface {
	ValidateAccount(ctx context.Context, username string, toAccountId int64, amount int64, currenct string) (*UserAccount, error)
	TransferTx(ctx context.Context, arg *TransferTxParams) (*UserAccount, error)
	GetBeneficiaries(ctx context.Context, username string) ([]*Beneficiary, error)
}

type service struct {
	repo Repo
}

func NewService(r Repo) (Service, error) {
	return &service{
		repo: r,
	}, nil
}

func (s *service) ValidateAccount(ctx context.Context, username string, toAccountId int64, amount int64, currency string) (*UserAccount, error) {

	//validate sender's account
	fromAccount, err := s.repo.GetAccountByUserName(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("account not found err")
			return nil, errors.ErrInvalidPayerAccount
		}
		fmt.Println("fromAccount err")
		return nil, err
	}

	if fromAccount.ID == toAccountId {
		return nil, errors.ErrInvalidBeneficiaryAccount
	}

	if fromAccount.Balance-amount < 0 {
		return nil, errors.ErrNoEnoughBalance
	}

	if fromAccount.Currency != currency {
		return nil, errors.ErrInvalidCurrency
	}

	//validate receiver's account
	toAccount, err := s.repo.GetAccountById(ctx, toAccountId)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("ToAccount not found err")
			return nil, errors.ErrInvalidBeneficiaryAccount
		}
		fmt.Println("ToAccount err")
		return nil, err
	}

	if fromAccount.Currency != toAccount.Currency {
		return nil, errors.ErrCurrencyMatch
	}

	return fromAccount, nil
}

func (s *service) TransferTx(ctx context.Context, arg *TransferTxParams) (*UserAccount, error) {
	result, err := s.repo.TransferTx(ctx, *arg)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) GetBeneficiaries(ctx context.Context, username string) ([]*Beneficiary, error) {
	result, err := s.repo.GetAllBeneficiaries(ctx, username)
	if err != nil {
		return nil, err
	}
	return result, nil
}
