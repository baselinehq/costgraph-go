package costgraph

import (
	"context"
	"errors"

	"github.com/baselinehq/costgraph-go/internal/api"
)

func (c *Client) ConnectAzureIntegration(ctx context.Context, request AzureConnectRequest) (*ConnectedIntegration, error) {
	response, err := c.api.ConnectAzureIntegration(ctx, request)
	result, err := decode[api.ConnectedIntegrationResponse]("connect Azure integration", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data.ID == "" {
		return nil, errors.New("connect Azure integration: missing connection UUID in CostGraph API response")
	}
	return &result.Data, nil
}

func (c *Client) GetAzureIntegrationSettings(ctx context.Context, id string) (*AzureConnectionSettings, error) {
	if err := requireConnectionID(id); err != nil {
		return nil, err
	}
	response, err := c.api.GetAzureIntegrationSettings(ctx, id)
	result, err := decode[api.AzureConnectionSettingsResponse]("read Azure integration settings", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data.ID == "" {
		return nil, errors.New("read Azure integration settings: missing connection in CostGraph API response")
	}
	return &result.Data, nil
}

func (c *Client) UpdateAzureIntegrationSettings(ctx context.Context, id string, request AzureConnectionUpdateRequest) (*AzureConnectionSettings, error) {
	if err := requireConnectionID(id); err != nil {
		return nil, err
	}
	response, err := c.api.UpdateAzureIntegrationSettings(ctx, id, request)
	result, err := decode[api.AzureConnectionSettingsResponse]("update Azure integration settings", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data.ID == "" {
		return nil, errors.New("update Azure integration settings: missing connection in CostGraph API response")
	}
	return &result.Data, nil
}
