package company

import "context"

type Interface interface {
	CreateCompany(ctx context.Context, name string) (*Company, error)
	GetCompany(ctx context.Context, name string) (*Company, error)
}

type Service struct {
	repo Repository
}

func (s Service) CreateCompany(ctx context.Context, name string) (*Company, error) {
	return s.repo.CreateCompany(ctx, name)
}
