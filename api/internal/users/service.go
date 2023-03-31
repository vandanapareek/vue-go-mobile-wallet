package users

import (
	"context"
	"fmt"
)

type Service interface {
	Login(ctx context.Context, username, password string) (*User, error)
	GetUserInfo(ctx context.Context, username string) (*UserInfo, error)
}

type service struct {
	repo Repo
}

func NewService(r Repo) (Service, error) {
	return &service{
		repo: r,
	}, nil
}

func (s *service) Login(ctx context.Context, username, password string) (*User, error) {
	user, err := s.repo.GetByUserName(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("Inexistant user")
	}
	if user.Password != password {
		return nil, fmt.Errorf("Wrong password")
	}
	return user, nil
}

func (s *service) GetUserInfo(ctx context.Context, username string) (*UserInfo, error) {
	user, err := s.repo.GetUserInfo(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("Inexistant user")
	}
	return user, nil
}
