package user

import (
	"backend/internal/model"
	"context"
	"errors"
)

type BasicUser struct {
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
}
type Interface interface {
	CreateUser(ctx context.Context, firstName, lastName, email, phoneNumber string) (string, error)
	CreateOrGetUser(ctx context.Context, firstName, lastName, email, phoneNumber string) (string, error)
}

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) CreateOrGetUser(ctx context.Context, firstName, lastName, email, phoneNumber string) (string, error) {
	u, err := s.repo.GetUserByEmail(ctx, email)
	if err == nil {
		return u.ID, nil
	}
	if !errors.Is(err, model.ErrNotFound) {
		return "", err
	}
	return s.CreateUser(ctx, firstName, lastName, email, phoneNumber)
}
func (s *Service) CreateUser(ctx context.Context, firstName, lastName, email, phoneNumber string) (string, error) {
	//TODO: validate args
	return s.repo.CreateUser(ctx, firstName, lastName, email, phoneNumber)
}
