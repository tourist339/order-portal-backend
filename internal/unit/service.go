package unit

import (
	"backend/internal/property"
	"backend/internal/roles"
	tenant2 "backend/internal/roles/tenant"
	"context"
	"fmt"
)

type Interface interface {
	CreateUnit(ctx context.Context, unitIdentifier string, propertyID string, tenants ...tenant2.Tenant) error
}

type Service struct {
	repo            Repository
	tenantService   roles.Interface
	propertyService property.Interface
}

func NewService(repo Repository, tenantService roles.Interface, propertyService property.Interface) *Service {
	return &Service{repo: repo, tenantService: tenantService, propertyService: propertyService}
}

func (s *Service) CreateUnit(ctx context.Context, propertyID, unitIdentifier string, tenants ...roles.TenantUser) error {
	//Check Property Exists
	_, err := s.propertyService.GetProperty(ctx, propertyID)
	if err != nil {
		fmt.Println("Error getting Property by", propertyID, err)
		return err
	}
	//Create Unit
	unitID, err := s.repo.CreateUnit(ctx, propertyID, unitIdentifier)
	if err != nil {
		fmt.Println("Error creating unit", err)
		return err
	}
	//TODO: Create Tenant in goroutines
	tenantIDs := []string{}
	for _, t := range tenants {
		tenantID, err := s.tenantService.CreateTenant(ctx, &t, unitID)
		if err != nil {
			fmt.Println("Error creating tenant", err)
			return err
		}
		tenantIDs = append(tenantIDs, tenantID)
	}
	// Attach Tenants to Unit
	if len(tenantIDs) > 0 {
		err = s.repo.AddTenantToUnit(ctx, tenantIDs, unitID)
		if err != nil {
			fmt.Println("Error adding tenants to unit", err)
			return err
		}
	}
	return nil
}
