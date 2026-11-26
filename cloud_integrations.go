package costgraph

import (
	"context"
	"errors"

	"github.com/baselinehq/costgraph-go/internal/api"
)

func (c *Client) ConnectAWSIntegration(ctx context.Context, request AWSConnectRequest) (*ConnectedIntegration, error) {
	response, err := c.api.ConnectAWSIntegration(ctx, request)
	result, err := decode[api.ConnectedIntegrationResponse]("connect AWS integration", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data.ID == "" {
		return nil, errors.New("connect AWS integration: missing connection UUID in CostGraph API response")
	}
	return &result.Data, nil
}

func (c *Client) GetAWSIntegrationSettings(ctx context.Context, id string) (*AWSConnectionSettings, error) {
	if err := requireConnectionID(id); err != nil {
		return nil, err
	}
	response, err := c.api.GetAWSIntegrationSettings(ctx, id)
	result, err := decode[api.AWSConnectionSettingsResponse]("read AWS integration settings", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data.ID == "" {
		return nil, errors.New("read AWS integration settings: missing connection in CostGraph API response")
	}
	return &result.Data, nil
}

func (c *Client) UpdateAWSIntegrationSettings(ctx context.Context, id string, request AWSConnectionUpdateRequest) (*AWSConnectionSettings, error) {
	if err := requireConnectionID(id); err != nil {
		return nil, err
	}
	response, err := c.api.UpdateAWSIntegrationSettings(ctx, id, request)
	result, err := decode[api.AWSConnectionSettingsResponse]("update AWS integration settings", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data.ID == "" {
		return nil, errors.New("update AWS integration settings: missing connection in CostGraph API response")
	}
	return &result.Data, nil
}

func (c *Client) ConnectGoogleCloudIntegration(ctx context.Context, request GCPConnectRequest) (*ConnectedIntegration, error) {
	response, err := c.api.ConnectGoogleCloudIntegration(ctx, request)
	result, err := decode[api.ConnectedIntegrationResponse]("connect Google Cloud integration", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data.ID == "" {
		return nil, errors.New("connect Google Cloud integration: missing connection UUID in CostGraph API response")
	}
	return &result.Data, nil
}

func (c *Client) GetGoogleCloudIntegrationSettings(ctx context.Context, id string) (*GCPConnectionSettings, error) {
	if err := requireConnectionID(id); err != nil {
		return nil, err
	}
	response, err := c.api.GetGoogleCloudIntegrationSettings(ctx, id)
	result, err := decode[api.GCPConnectionSettingsResponse]("read Google Cloud integration settings", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data.ID == "" {
		return nil, errors.New("read Google Cloud integration settings: missing connection in CostGraph API response")
	}
	return &result.Data, nil
}

func (c *Client) UpdateGoogleCloudIntegrationSettings(ctx context.Context, id string, request GCPConnectionUpdateRequest) (*GCPConnectionSettings, error) {
	if err := requireConnectionID(id); err != nil {
		return nil, err
	}
	response, err := c.api.UpdateGoogleCloudIntegrationSettings(ctx, id, request)
	result, err := decode[api.GCPConnectionSettingsResponse]("update Google Cloud integration settings", response, err)
	if err != nil {
		return nil, err
	}
	if result.Data.ID == "" {
		return nil, errors.New("update Google Cloud integration settings: missing connection in CostGraph API response")
	}
	return &result.Data, nil
}
