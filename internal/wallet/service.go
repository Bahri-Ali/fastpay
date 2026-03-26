package wallet

import (
	"context"
	"errors"
)

type Service interface {
	GetWallet(ctx context.Context, userID string) (*Wallet, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetWallet(ctx context.Context, userID string) (*Wallet, error) {
	wallet, err := s.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if wallet == nil {
		return nil, errors.New("wallet not found")
	}

	return wallet, nil
}
