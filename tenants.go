package costgraph

import (
	"context"
	"errors"

	"github.com/baselinehq/costgraph-go/internal/api"
	"github.com/google/uuid"
)

const ScopeOrganizationManage = "organization:manage"

func (c *Client) ListTenants(ctx context.Context, organizationID string) ([]ManagedTenant, error) {
	if err := requireTenantIDs(organizationID); err != nil {
		return nil, err
	}
	response, err := c.api.ListOrganizationTenants(ctx, organizationID)
	result, err := decode[api.ManagedTenantListResponse]("list organization tenants", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data == nil {
		return nil, errors.New("list organization tenants: missing tenant list in CostGraph API response")
	}
	return result.Data, nil
}

func (c *Client) GetTenant(ctx context.Context, organizationID, tenantID string) (*ManagedTenant, error) {
	if err := requireTenantIDs(organizationID, tenantID); err != nil {
		return nil, err
	}
	response, err := c.api.GetOrganizationTenant(ctx, organizationID, tenantID)
	result, err := decode[api.ManagedTenantResponse]("read organization tenant", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data.ID == "" {
		return nil, errors.New("read organization tenant: missing tenant in CostGraph API response")
	}
	return &result.Data, nil
}

func (c *Client) CreateTenant(ctx context.Context, organizationID string, request CreateTenantRequest) (*ManagedTenant, error) {
	if err := requireTenantIDs(organizationID); err != nil {
		return nil, err
	}
	response, err := c.api.CreateOrganizationTenant(ctx, organizationID, request)
	result, err := decode[api.ManagedTenantResponse]("create organization tenant", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data.ID == "" {
		return nil, errors.New("create organization tenant: missing tenant in CostGraph API response")
	}
	return &result.Data, nil
}

func (c *Client) UpdateTenant(ctx context.Context, organizationID, tenantID string, request UpdateTenantRequest) (*ManagedTenant, error) {
	if err := requireTenantIDs(organizationID, tenantID); err != nil {
		return nil, err
	}
	response, err := c.api.UpdateOrganizationTenant(ctx, organizationID, tenantID, request)
	result, err := decode[api.ManagedTenantResponse]("update organization tenant", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data.ID == "" {
		return nil, errors.New("update organization tenant: missing tenant in CostGraph API response")
	}
	return &result.Data, nil
}

func (c *Client) DeleteTenant(ctx context.Context, organizationID, tenantID string) error {
	if err := requireTenantIDs(organizationID, tenantID); err != nil {
		return err
	}
	response, err := c.api.DeleteOrganizationTenant(ctx, organizationID, tenantID)
	_, err = decode[api.SuccessResponse]("delete organization tenant", response, err)
	return err
}

func requireTenantIDs(organizationID string, tenantIDs ...string) error {
	if _, err := uuid.Parse(organizationID); err != nil {
		return errors.New("organization id must be a UUID")
	}
	for _, tenantID := range tenantIDs {
		if _, err := uuid.Parse(tenantID); err != nil {
			return errors.New("tenant id must be a UUID")
		}
	}
	return nil
}
